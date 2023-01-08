package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// type Setting struct {
// 	AutoReply      bool   `json:"auto_reply"`
// 	AutoReplyGroup bool   `json:"auto_reply_group"`
// 	AutoBot        string `json:"auto_bot"`
// 	AutoDesc       string `json:"auto_desc"`
// }

type Store struct {
	Data        []map[string]interface{}
	CurrentData map[string]interface{}
}

/** 初始化 stroe */
func (s *Store) InitStroe() {
	data, err := ioutil.ReadFile(CONFIG_PATH)
	if err != nil {
		fmt.Println("读取配置文件失败", err)
		return
	}
	if len(data) == 0 {
		return
	}
	err1 := json.Unmarshal(data, &s.Data)
	if err1 != nil {
		fmt.Println("配置解析失败", err)
		return
	}
	fmt.Println("配置解析===", s.Data)
}

func (s *Store) Save(config []interface{}) {
	file, err := os.Create(CONFIG_PATH)
	if err != nil {
		fmt.Println("保存配置路径读取失败", err)
		return
	}
	defer file.Close()
	newConfigs := append(append(s.Data[:0], config[0].(map[string]interface{})), s.Data[0:]...)

	n := len(newConfigs)
	var resConfig []map[string]interface{}
	cache := make(map[string]int)
	for i := 0; i < n; i++ {
		forFriend := ""
		if newConfigs[i]["for_friend"] != nil {
			forFriend = newConfigs[i]["for_friend"].(string)
		}
		unid := ""
		if newConfigs[i]["unid"] != nil {
			unid = newConfigs[i]["unid"].(string)
		}
		cacheKey := unid + forFriend
		if cache[cacheKey] == 0 {
			resConfig = append(resConfig, newConfigs[i])
		}
		cache[cacheKey] = 1
	}

	//将map类型的i转换成[]byte类型
	data, err1 := json.MarshalIndent(resConfig, "", " ")
	if err1 != nil {
		fmt.Println("保存配置转换数据失败", err)
		return
	}

	// //将数据data写入文件filePath中，并且修改文件权限为755
	// if err = ioutil.WriteFile(fileinfo, data, 0755)

	file.Write(data)
}

func (s *Store) Update(unid string) {
	s.InitStroe()
	s.GetDataById(unid)
}

func (s *Store) GetDataById(unid string) interface{} {
	for i := 0; i < len(s.Data); i++ {
		if s.Data[i]["unid"] == unid {
			s.CurrentData = s.Data[i]
			fmt.Println("GetDataById=== s.Data[i] ", s.CurrentData)
			return s.Data[i]
		}
	}
	return nil
}
