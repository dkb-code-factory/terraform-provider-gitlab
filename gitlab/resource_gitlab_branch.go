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
		},
	}
}

func resourceGitlabBranchSetToState(d *schema.ResourceData, branch *gitlab.Branch) {
	d.SetId(branch.Name)
	d.Set("name", branch.Name)
	d.Set("protected", branch.Protected)
	d.Set("merged", branch.Merged)
	d.Set("default", branch.Default)
	d.Set("developers_can_push", branch.DevelopersCanPush)
	d.Set("developers_can_merge", branch.DevelopersCanMerge)
}

func resourceGitlabBranchCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.CreateBranchOptions{
		Branch: gitlab.String(d.Get("name").(string)),
		Ref:    gitlab.String(d.Get("reference_branch").(string)),
	}

	log.Printf("[DEBUG] create gitlab branch %s", *options.Branch)

	branch, _, err := client.Branches.CreateBranch(project, options)
	if err != nil {
		return err
	}

	d.SetId(branch.Name)

	return resourceGitlabBranchRead(d, meta)
}

func resourceGitlabBranchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branchName := d.Get("name").(string)
	log.Printf("[DEBUG] read gitlab branch %s/%s", project, branchName)

	branch, _, err := client.Branches.GetBranch(project, branchName)
	if err != nil {
		return err
	}

	resourceGitlabBranchSetToState(d, branch)
	return nil
}

func resourceGitlabBranchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branch := d.Get("name").(string)

	log.Printf("[DEBUG] Delete branch %s for project %s", branch, project)

	_, err := client.Branches.DeleteBranch(project, branch)
	return err
}
