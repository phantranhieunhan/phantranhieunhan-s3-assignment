package payload

import "github.com/phantranhieunhan/s3-assignment/pkg/util"

type SubscriberUserPayload struct {
	Requestor string
	Target    string
}

type SubscriberUserPayloads []SubscriberUserPayload

func (s SubscriberUserPayloads) GetEmails(total int) []string {
	userIds := make([]string, 0, len(s)*total)
	for _, u := range s {
		if !util.IsContain(userIds, u.Requestor) {
			userIds = append(userIds, u.Requestor)
		}
		if !util.IsContain(userIds, u.Target) {
			userIds = append(userIds, u.Target)
		}
	}

	return userIds
}
