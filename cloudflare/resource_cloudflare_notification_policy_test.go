package cloudflare

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloudflareNotificationPolicy_Basic(t *testing.T) {
	rnd := generateRandomResourceName()
	resourceName := "cloudflare_notification_policy." + rnd
	updatedPolicyName := "updated test SSL policy from terraform provider"
	updatedPolicyDesc := "updated description"
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckAccount(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCheckCloudflareNotificationPolicy(rnd, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test SSL policy from terraform provider"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alert_type", "universal_ssl_event_type"),
					resource.TestCheckResourceAttr(resourceName, "account_id", accountID),
				),
			},
			{
				Config: testCheckCloudflareNotificationPolicyUpdated(rnd, updatedPolicyName, updatedPolicyDesc, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedPolicyName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedPolicyDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alert_type", "universal_ssl_event_type"),
					resource.TestCheckResourceAttr(resourceName, "account_id", accountID),
				),
			},
		},
	})
}

func testCheckCloudflareNotificationPolicy(name, accountID string) string {
	return fmt.Sprintf(`
  resource "cloudflare_notification_policy" "%[1]s" {
    name        = "test SSL policy from terraform provider"
    account_id  = "%[2]s"
    description = "test description"
    enabled     =  true
    alert_type  = "universal_ssl_event_type"
    email_integration {
      name =  ""
      id   =  "test@example.com"
    }
  }`, name, accountID)
}

func testCheckCloudflareNotificationPolicyUpdated(resName, policyName, policyDesc, accountID string) string {
	return fmt.Sprintf(`
  resource "cloudflare_notification_policy" "%[1]s" {
    name        = "%[2]s"
    account_id  = "%[4]s"
    description = "%[3]s"
    enabled     =  true
    alert_type  = "universal_ssl_event_type"
    email_integration {
      name =  ""
      id   =  "test@example.com"
    }
  }`, resName, policyName, policyDesc, accountID)
}

func TestAccCloudflareNotificationPolicy_WithFiltersAttribute(t *testing.T) {
	rnd := generateRandomResourceName()
	resourceName := "cloudflare_notification_policy." + rnd
	updatedPolicyName := "updated workers usage notification"
	updatedPolicyDesc := "updated description"
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckAccount(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testCheckCloudflareNotificationPolicyWithFiltersAttribute(rnd, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "workers usage notification"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alert_type", "billing_usage_alert"),
					resource.TestCheckResourceAttr(resourceName, "account_id", accountID),
					resource.TestCheckTypeSetElemAttr(resourceName, "filters.0.product.*", "worker_requests"),
					resource.TestCheckTypeSetElemAttr(resourceName, "filters.0.limit.*", "100"),
				),
			},
			{
				Config: testCheckCloudflareNotificationPolicyWithFiltersAttributeUpdated(rnd, updatedPolicyName, updatedPolicyDesc, accountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedPolicyName),
					resource.TestCheckResourceAttr(resourceName, "description", updatedPolicyDesc),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alert_type", "billing_usage_alert"),
					resource.TestCheckResourceAttr(resourceName, "account_id", accountID),
					resource.TestCheckTypeSetElemAttr(resourceName, "filters.0.product.*", "worker_requests"),
					resource.TestCheckTypeSetElemAttr(resourceName, "filters.0.limit.*", "100"),
				),
			},
		},
	})
}

func testCheckCloudflareNotificationPolicyWithFiltersAttribute(name, accountID string) string {
	return fmt.Sprintf(`
  resource "cloudflare_notification_policy" "%[1]s" {
    name        = "workers usage notification"
    account_id  = "%[2]s"
    description = "test description"
    enabled     =  true
    alert_type  = "billing_usage_alert"
    email_integration {
      name =  ""
      id   =  "test@example.com"
    }
    filters {
      product = [
        "worker_requests",
      ]
	  limit = ["100"]
	}
  }`, name, accountID)
}

func testCheckCloudflareNotificationPolicyWithFiltersAttributeUpdated(name, policyName, policyDesc, accountID string) string {
	return fmt.Sprintf(`
  resource "cloudflare_notification_policy" "%[1]s" {
    name        = "%[2]s"
    account_id  = "%[4]s"
    description = "%[3]s"
    enabled     =  true
    alert_type  = "billing_usage_alert"
    email_integration {
      name =  ""
      id   =  "test@example.com"
    }
    filters {
      product = [
        "worker_requests",
      ]
      limit = ["100"]
	}
  }`, name, policyName, policyDesc, accountID)
}

func TestFlattenExpandFilters(t *testing.T) {
	filters := map[string][]string{
		"services": []string{"waf", "firewallrules"},
		"zones":    []string{"abc123"},
	}
	flattenedFilters := flattenNotificationPolicyFilter(filters)
	expandedFilters := expandNotificationPolicyFilter(flattenedFilters)
	for k, _ := range filters {
		sort.Strings(filters[k])
		sort.Strings(expandedFilters[k])
		assert.EqualValuesf(t, filters[k], expandedFilters[k], "values should equal without order")
	}
}
