package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/poseidon/matchbox/matchbox/storage/storagepb"
	"github.com/sirupsen/logrus"
)

type etcdStore struct {
	config clientv3.Config
	requestTimeout time.Duration

	logger *logrus.Logger
}

func NewEtcdStore(endpoints []string, logger *logrus.Logger) Store {
	return &etcdStore{
		config: clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		},
		logger: logger,
		requestTimeout: 1 * time.Second,
	}
}

func (e *etcdStore) putKeyValue(key string, value []byte) (err error) {
	e.logger.WithFields(logrus.Fields{"function": "putKeyValue", "value": value}).Debugf("Putting key %s")

	cli, err := clientv3.New(e.config)
	if err == nil {
		defer cli.Close()
		ctx, cancel := e.createContextAndCancelForEtcdClient()
		_, err = cli.Put(ctx, key, string(value))
		cancel()
	}

	if err != nil {
		e.logger.WithFields(logrus.Fields{"function": "putKeyValue", "value": value, "key": key}).Error(err)
	}

	return err
}

func (e *etcdStore) createContextAndCancelForEtcdClient() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), e.requestTimeout)
}

func (e *etcdStore) getKey(key string) (result []byte, err error) {
	e.logger.WithFields(logrus.Fields{"function": "getKey"}).Debugf("Getting key %s", key)

	cli, err := clientv3.New(e.config)
	if err == nil {
		var resp *clientv3.GetResponse
		defer cli.Close()
		ctx, cancel := e.createContextAndCancelForEtcdClient()

		resp, err = cli.Get(ctx, key)
		cancel()

		for _, ev := range resp.Kvs {
			result = ev.Value
		}
	}

	if err != nil {
		e.logger.WithFields(logrus.Fields{"function": "getKey", "key": key}).Error(err)
	}

	e.logger.WithFields(logrus.Fields{"function": "getKey"}).Debugf("Result: %s", result)
	return result, err
}

func (e *etcdStore) getPrefixedKeys(key string) (result [][]byte, err error) {
	e.logger.WithFields(logrus.Fields{"function": "getPrefixedKeys"}).Debugf("Getting key %s", key)

	cli, err := clientv3.New(e.config)
	if err == nil {
		var resp *clientv3.GetResponse
		defer cli.Close()
		ctx, cancel := e.createContextAndCancelForEtcdClient()
		resp, err = cli.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
		cancel()

		for _, ev := range resp.Kvs {
			result = append(result, ev.Key)
			e.logger.WithFields(logrus.Fields{"function": "getPrefixedKeys"}).Debugf("Result: %s", ev.Key)
		}
	}

	if err != nil {
		e.logger.WithFields(logrus.Fields{"function": "getPrefixedKeys", "key": key}).Error(err)
	}

	return result, err
}

func (e *etcdStore) deleteKey(key string) (err error) {
	e.logger.WithFields(logrus.Fields{"function": "deleteKey"}).Debugf("Deleting key %s", key)

	cli, err := clientv3.New(e.config)
	if err == nil {
		defer cli.Close()
		ctx, cancel := e.createContextAndCancelForEtcdClient()
		_, err = cli.Delete(ctx, key)
		cancel()
	}

	if err != nil {
		e.logger.WithFields(logrus.Fields{"function": "deleteKey", "key": key}).Error(err)
	}

	return err
}

func (e *etcdStore) GroupPut(group *storagepb.Group) (err error) {
	e.logger.WithFields(logrus.Fields{"function": "GroupPut"}).Debugf("Putting group %s", group)

	var key string
	var data []byte

	richGroup, err := group.ToRichGroup()
	if err == nil {
		data, err = e.createPrettyJSONFromStruct(richGroup)
		if err == nil {
			key = e.createGroupKeyFromId(group.Id)
			err = e.putKeyValue(key, data)
		}
	}

	if err != nil {
		e.logger.WithFields(logrus.Fields{"function": "GroupPut", "group": group}).Error(err)
	}

	return err
}

func (e *etcdStore) createPrettyJSONFromStruct(data any) ([]byte, error) {
	return json.MarshalIndent(data, "", "\t")
}

func (e *etcdStore) GroupGet(id string) (result *storagepb.Group, err error) {
	e.logger.WithFields(logrus.Fields{"function": "GroupGet"}).Debugf("Getting the group %s", id)

	key := e.createGroupKeyFromId(id)
	data, err := e.getKey(key)
	if err == nil {
		result, err = storagepb.ParseGroup(data)
	}

	if err != nil {
		e.logger.WithFields(logrus.Fields{"function": "GroupGet", "id": id}).Error(err)
	}

	return result, err
}

func (e *etcdStore) createGroupKeyFromId(id string) string {
	return fmt.Sprintf("groups/%s", id)
}

func (e *etcdStore) GroupDelete(id string) error {
	e.logger.WithFields(logrus.Fields{"function": "GroupDelete"}).Debugf("Deleting the group %s", id)

	key := fmt.Sprintf("groups/%s", id)
	return e.deleteKey(key)
}

func (e *etcdStore) GroupList() (groups []*storagepb.Group, err error) {
	e.logger.WithFields(logrus.Fields{"function": "GroupList"}).Debug("Getting the groups")

	keys, err := e.getPrefixedKeys("groups/")
	if err == nil {
		e.logger.WithFields(logrus.Fields{"function": "GroupList"}).Debugf("Found %v groups", len(keys))
		groups = make([]*storagepb.Group, 0, len(keys))

		for _, key := range keys {
			name := e.getGroupNameFromKey(key)
			group, err := e.GroupGet(name)
			if err == nil {
				groups = append(groups, group)
			}
		}
	}

	return groups, err
}

func (e *etcdStore) getGroupNameFromKey(key []byte) string {
	return strings.TrimPrefix(string(key), "groups/")
}

func (e *etcdStore) ProfilePut(profile *storagepb.Profile) (err error) {
	e.logger.WithFields(logrus.Fields{"function": "ProfilePut"}).Debugf("Putting the profile %s", profile)

	data, err := json.MarshalIndent(profile, "", "\t")
	if err == nil {
		key := fmt.Sprintf("profiles/%s", profile.Id)
		err = e.putKeyValue(key,data)
	}

	return err
}

func (e *etcdStore) ProfileGet(id string) (profile *storagepb.Profile, err error) {
	e.logger.WithFields(logrus.Fields{"function": "ProfileGet"}).Debugf("Getting the profile %s", id)

	key := fmt.Sprintf("profiles/%s", id)
	data, err := e.getKey(key)
	if err == nil {
		profile, err = e.createProfileFromJSONRawData(data)
	}

	e.logger.WithFields(logrus.Fields{"function": "ProfileGet"}).Debugf("Result: %s", profile)
	return profile, err
}

func (e *etcdStore) createProfileFromJSONRawData(data []byte) (profile *storagepb.Profile, err error) {
	profile = new(storagepb.Profile)
	err = json.Unmarshal(data, profile)
	if err == nil {
		err = profile.AssertValid()
	} else {
		e.logger.WithFields(logrus.Fields{"function": "createProfileFromJSONRawData"}).Error(err)
	}

	return profile, err
}

func (e *etcdStore) ProfileDelete(id string) error {
	e.logger.WithFields(logrus.Fields{"function": "ProfileDelete"}).Debugf("Deleting the profile %s", id)

	key := fmt.Sprintf("profiles/%s", id)
	return e.deleteKey(key)
}

func (e *etcdStore) ProfileList() (profiles []*storagepb.Profile, err error) {
	e.logger.WithFields(logrus.Fields{"function": "ProfileList"}).Debug("Listing the profiles")

	keys, err := e.getPrefixedKeys("profiles/")
	if err == nil {
		profiles = make([]*storagepb.Profile, 0, len(keys))
		for _, key := range keys {
			name := strings.TrimPrefix(string(key), "profiles/")
			profile, err := e.ProfileGet(name)
			if err == nil {
				profiles = append(profiles, profile)
			}
		}
	}

	e.logger.WithFields(logrus.Fields{"function": "ProfileList"}).Debugf("Result: %s", profiles)
	return profiles, err
}

func (e *etcdStore) IgnitionPut(name string, config []byte) error {
	e.logger.WithFields(logrus.Fields{"function": "IgnitionPut", "value": config}).Debugf("Putting ignition %s", name)

	key := fmt.Sprintf("ignitions/%s", name)
	return e.putKeyValue(key, config)
}

func (e *etcdStore) IgnitionGet(name string) (string, error) {
	e.logger.WithFields(logrus.Fields{"function": "IgnitionGet"}).Debugf("Getting ignition %s", name)

	key := fmt.Sprintf("ignitions/%s", name)
	data, err := e.getKey(key)
	return string(data), err
}

func (e *etcdStore) IgnitionDelete(name string) error {
	e.logger.WithFields(logrus.Fields{"function": "IgnitionDelete"}).Debugf("Deleting ignition %s", name)

	key := fmt.Sprintf("ignitions/%s", name)
	return e.deleteKey(key)
}

func (e *etcdStore) GenericPut(name string, config []byte) error {
	e.logger.WithFields(logrus.Fields{"function": "GenericPut"}).Debugf("Putting generic %s", name)

	key := fmt.Sprintf("generics/%s", name)
	return e.putKeyValue(key, config)
}

func (e *etcdStore) GenericGet(name string) (string, error) {
	e.logger.WithFields(logrus.Fields{"function": "GenericGet"}).Debugf("Getting generic %s", name)

	key := fmt.Sprintf("generics/%s", name)
	data, err := e.getKey(key)
	return string(data), err
}

func (e *etcdStore) GenericDelete(name string) error {
	e.logger.WithFields(logrus.Fields{"function": "GenericDelete"}).Debugf("Deleting generic %s", name)

	key := fmt.Sprintf("generics/%s", name)
	return e.deleteKey(key)
}

func (e *etcdStore) CloudGet(name string) (string, error) {
	e.logger.WithFields(logrus.Fields{"function": "CloudGet"}).Debugf("Getting cloud %s", name)

	key := fmt.Sprintf("clouds/%s", name)
	data, err := e.getKey(key)
	return string(data), err
}