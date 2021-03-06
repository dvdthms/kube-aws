package model

import (
	"errors"
	"fmt"
	"regexp"
)

type IAMConfig struct {
	Role            IAMRole            `yaml:"role,omitempty"`
	InstanceProfile IAMInstanceProfile `yaml:"instanceProfile,omitempty"`
	UnknownKeys     `yaml:",inline"`
}

type IAMRole struct {
	ARN             `yaml:",inline"`
	Name            string             `yaml:"name,omitempty"`
	ManagedPolicies []IAMManagedPolicy `yaml:"managedPolicies,omitempty"`
}

type IAMManagedPolicy struct {
	ARN `yaml:",inline"`
}

type IAMInstanceProfile struct {
	ARN `yaml:",inline"`
}

func (c IAMConfig) Validate() error {
	if c.InstanceProfile.Arn != "" && c.Role.Name != "" {
		return errors.New("failed to parse `iam` config: either you set `role.*` options or `instanceProfile.arn` ones but not both")
	}
	if c.InstanceProfile.Arn != "" && len(c.Role.ManagedPolicies) > 0 {
		return errors.New("failed to parse `iam` config: either you set `role.*` options or `instanceProfile.arn` ones but not both")
	}

	managedPolicyRegexp := regexp.MustCompile(`arn:aws:iam::((\d{12})|aws):policy/([a-zA-Z0-9-=,\\.@_]{1,128})`)
	instanceProfileRegexp := regexp.MustCompile(`arn:aws:iam::(\d{12}):instance-profile/([a-zA-Z0-9-=,\\.@_]{1,128})`)
	for _, policy := range c.Role.ManagedPolicies {
		if !managedPolicyRegexp.MatchString(policy.Arn) {
			return fmt.Errorf("invalid managed policy arn, your managed policy must match this (=arn:aws:iam::(YOURACCOUNTID|aws):policy/POLICYNAME), provided this (%s)", policy.Arn)
		}
	}
	if c.InstanceProfile.Arn != "" {
		if !instanceProfileRegexp.MatchString(c.InstanceProfile.Arn) {
			return fmt.Errorf("invalid instance profile, your instance profile must match (=arn:aws:iam::YOURACCOUNTID:instance-profile/INSTANCEPROFILENAME), provided (%s)", c.InstanceProfile.Arn)
		}

	}
	return nil

}
