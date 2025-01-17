/*
Copyright (c) 2022 PaddlePaddle Authors. All Rights Reserve.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sortpolicy

import (
	"paddleflow/pkg/job/api"
)

// PolicyName indicates name of queue sort policy.
const PriorityPolicyName = "priority"

type priorityPolicy struct {
	// Arguments given for the sort policy
	policyArguments api.Arguments
}

// New return priority sort policy
func PriorityPolicyNew(arguments api.Arguments) (api.SortPolicy, error) {
	return &priorityPolicy{
		policyArguments: arguments,
	}, nil
}

func (pp *priorityPolicy) Name() string {
	return PriorityPolicyName
}

func (pp *priorityPolicy) OrderFn(l, r interface{}) int {
	lv := l.(*api.PFJob)
	rv := r.(*api.PFJob)

	if lv.Priority > rv.Priority {
		return -1
	}

	if lv.Priority < rv.Priority {
		return 1
	}

	return 0
}
