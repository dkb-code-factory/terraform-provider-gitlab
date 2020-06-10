package gitlab

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func ResourceGitlabBranch() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabBranchCreate,
		Read:   resourceGitlabBranchRead,
		Delete: resourceGitlabBranchDelete,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"reference_branch": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"merge_access_level": {
				Type:         schema.TypeString,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
				Required:     true,
				ForceNew:     true,
			},
			"push_access_level": {
				Type:         schema.TypeString,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
				Required:     true,
				ForceNew:     true,
			},
		},
	}
}

func resourceGitlabBranchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.CreateBranchOptions{
		Branch: gitlab.String(d.Get("name").(string)),
		Ref:    gitlab.String(d.Get("reference_branch").(string)),
	}

	log.Printf("[DEBUG] create gitlab branch %s", *options.Name)

	branch, _, err := clinet.Branches.CreateBranch(project, options)
	if err != nil {
		return err
	}

	d.SetId(branch.Name)

	return resourceGitlabBranchRead(d, meta)
}

func resourceGitlabBranchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branchName := d.Id()
	log.Printf("[DEBUG] read gitlab branch %s/%s", project, branchName)

	page := 1
	labelsLen := 0
	for page == 1 || labelsLen != 0 {
		braches, _, err := client.Branches.ListBranches(project, &gitlab.ListBranchesOptions{Page: page})
		if err != nil {
			return err
		}
		for _, branch := range branches {
			if branch.Name == branchName {
				d.Set("name", branch.Name)
				return nil
			}
		}
		branchesLen = len(branches)
		page = page + 1
	}

	log.Printf("[DEBUG] failed to read gitlab branch %s/%s", project, branchName)
	d.SetId("")
	return nil
}

func resourceGitlabBranchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branch := d.Get("name").(string)

	log.Printf("[DEBUG] Delete branch %s for project %s", branch, project)

	_, err := client.Branches.DeleteBranch(branch)
	return err
}
