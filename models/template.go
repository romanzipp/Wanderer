package models

import (
	"errors"
	"fmt"
	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Template struct {
	gorm.Model
	Name       string
	NomadJobID string
	Content    string
	Server     Server
	ServerID   uint
	Versions   []TemplateVersion
	Locked     bool
}

func (t Template) GetNomadJobUrl() string {
	return fmt.Sprintf("%s/ui/jobs//%s", t.Server.Address, t.NomadJobID)
}

func (t Template) getPopulatedContentForVersion(db *gorm.DB, newVersion TemplateVersion) string {
	// get all previous versions from database
	var versions []TemplateVersion
	db.Where("id != ?", newVersion.ID).Where("template_id = ?", t.ID).Find(&versions)

	content := t.Content

	// add new version to slice
	versions = append(versions, newVersion)

	// replace selectors with version strings
	for _, version := range versions {
		content = strings.Replace(content, fmt.Sprintf("{{ %s }}", version.Selector), version.LastVersion, -1)
		log.Debug().Msgf("version replacement: %s -> %s", version.Selector, version.LastVersion)
	}

	return content
}

func (t Template) Deploy(db *gorm.DB, templateVersion *TemplateVersion, versionString string) error {
	if t.Server.ID == 0 {
		return errors.New("missing server")
	}

	// get new template content with new versions
	payload := t.getPopulatedContentForVersion(db, *templateVersion)

	err := t.deployPayload(db, payload)
	if err != nil {
		return err
	}

	templateVersion.LastDeployedAt = time.Now()

	db.Save(templateVersion)

	return nil
}

// --------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------

func (t Template) getPopulatedContent(db *gorm.DB) string {
	// get all versions from database
	var versions []TemplateVersion
	db.Where("template_id = ?", t.ID).Find(&versions)

	content := t.Content

	// replace selectors with version strings
	for _, version := range versions {
		content = strings.Replace(content, fmt.Sprintf("{{ %s }}", version.Selector), version.LastVersion, -1)
		log.Debug().Msgf("version replacement: %s -> %s", version.Selector, version.LastVersion)
	}

	return content
}

func (t Template) DeployCurrent(db *gorm.DB) error {
	if t.Server.ID == 0 {
		return errors.New("missing server")
	}

	// get new template content with new versions
	payload := t.getPopulatedContent(db)

	err := t.deployPayload(db, payload)
	if err != nil {
		return err
	}

	return nil
}

// --------------------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------

func (t Template) deployPayload(db *gorm.DB, payload string) error {
	if t.Server.ID == 0 {
		return errors.New("missing server")
	}

	client, err := t.Server.NewNomadClient()
	if err != nil {
		return err
	}

	// parse job spec
	job, err := client.Jobs().ParseHCL(payload, false)
	if err != nil {
		log.Debug().Err(err)
		return err
	}

	// plan job with new template
	plan, _, err := client.Jobs().Plan(job, false, &api.WriteOptions{})
	if err != nil {
		log.Debug().Err(err)
		return err
	}

	log.Debug().Msgf("planned job '%s' mit modify index %s", *job.ID, plan.JobModifyIndex)

	// dispatch job
	res, _, err := client.Jobs().Register(job, &api.WriteOptions{})
	if err != nil {
		log.Debug().Err(err)
		return err
	}

	log.Debug().Msgf("dispatched job '%s' mit modify index %s", *job.ID, res.JobModifyIndex)

	return nil
}
