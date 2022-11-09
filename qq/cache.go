package qq

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/common-utils/request"
)

const (
	cacheDirPath = "cache/"
	imagePath    = "images/"
	essencePath  = "essences/"
)

// images

func saveGroupImages(msg *message.GroupMessage) {
	err := os.MkdirAll(cacheDirPath+imagePath, os.ModePerm)
	if err != nil {
		logger.Errorf("創建群图片緩存資料夾時出現錯誤: %v", err)
		return
	}

	for _, element := range msg.Elements {

		var url string
		var imageId string
		var hash []byte

		switch e := element.(type) {
		case *message.FriendImageElement:
			imageId, hash, url = e.ImageId, e.Md5, e.Url
		case *message.GroupImageElement:
			if e.Flash || e.Url == "" {
				if url, err := bot.Instance.GetGroupImageDownloadUrl(e.FileId, msg.GroupCode, e.Md5); err == nil {
					e.Url = url
				} else {
					logger.Errorf("圖片URL為空或是閃照, 但嘗試獲取圖片 %s 的下載URL時出現錯誤: %v", e.FileId, err)
				}
			}

			imageId, hash, url = e.ImageId, e.Md5, e.Url
		case *message.GuildImageElement:
			imageId, hash, url = fmt.Sprint(e.FileId), e.Md5, e.Url
		default:
			continue
		}

		name := hex.EncodeToString(hash)

		b, err := request.GetBytesByUrl(url)
		if err != nil {
			logger.Errorf("下載圖片 %s 時出現錯誤: %v", strings.ToLower(imageId), name, err)
			continue
		}
		err = os.WriteFile(cacheDirPath+imagePath+name, b, os.ModePerm)
		if err != nil {
			logger.Errorf("緩存圖片 %s 時出現錯誤: %v", strings.ToLower(imageId), err)
		} else {
			logger.Infof("緩存圖片 %s 成功。", strings.ToLower(imageId))
		}
	}
}

func fixGroupImages(gp int64, sending *message.GroupMessage) {
	fixed := make([]message.IMessageElement, len(sending.Elements))
	for _, element := range sending.Elements {
		if groupImage, ok := element.(*message.GroupImageElement); ok {
			name := hex.EncodeToString(groupImage.Md5)
			b, err := os.ReadFile(cacheDirPath + "images/" + name)

			var img *message.GroupImageElement

			if err == nil {
				img, err = NewImagesByByteWithGroup(gp, b)
				if err != nil {
					logger.Errorf("群圖片上傳失敗: %v, 將使用QQ查詢", err)
				} else {
					logger.Infof("恢复缓存图片 %s 成功。", strings.ToLower(groupImage.ImageId))
				}
			} else {

				logger.Errorf("讀取緩存文件 %s 時出現錯誤: %v, 將使用QQ查詢", name, err)

				if url, err := bot.Instance.GetGroupImageDownloadUrl(groupImage.FileId, gp, groupImage.Md5); err == nil {
					logger.Infof("获取群图片下载链接成功，将尝试使用上传通道")
					img, err = NewImageByUrlWithGroup(gp, url)
					if err == nil {
						logger.Infof("群图片上传成功")
					} else {
						logger.Warnf("群图片上传失败: %v", err)
					}
				}
			}

			if img == nil {
				img, err = bot.Instance.QueryGroupImage(gp, groupImage.Md5, groupImage.Size)
				if err != nil {
					logger.Errorf("QQ查詢群圖片失敗: %v, 將繼續使用舊元素發送。", err)
					img = groupImage
				} else {
					logger.Infof("查询图片 %s 成功。", strings.ToLower(groupImage.ImageId))

					//查詢成功后下載
					url := img.Url
					b, err := request.GetBytesByUrl(url)
					if err != nil {
						logger.Errorf("下載查詢圖片 %s 時出現錯誤: %v", strings.ToLower(groupImage.ImageId), name, err)
					} else {
						err = os.WriteFile(cacheDirPath+imagePath+name, b, os.ModePerm)
						if err != nil {
							logger.Errorf("緩存查詢圖片 %s 時出現錯誤: %v", strings.ToLower(groupImage.ImageId), err)
						} else {
							logger.Infof("緩存查詢圖片 %s 成功。", strings.ToLower(groupImage.ImageId))
						}
					}
				}
			}

			// ensure not null
			if img == nil {
				logger.Warn("檢測到圖片為 null, 將繼續使用舊元素發送。")
				img = groupImage
			}

			fixed = append(fixed, img)
		} else {
			fixed = append(fixed, element)
		}
	}

	sending.Elements = fixed
}

// essence

func saveGroupEssence(msg *message.GroupMessage) {
	_ = saveGroupEssenceErr(msg)
}

func saveGroupEssenceErr(msg *message.GroupMessage) error {
	err := os.MkdirAll(cacheDirPath+essencePath, os.ModePerm)
	if err != nil {
		logger.Errorf("創建群精华緩存資料夾時出現錯誤: %v", err)
		return err
	}

	persit := &PersistentGroupMessage{}
	err = persit.Parse(msg)
	if err != nil {
		logger.Errorf("尝试持久化群精华消息时出现错误: %v", err)
		return err
	}

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(persit)
	if err != nil {
		logger.Errorf("尝试序列化群精华消息时出现错误: %v", err)
		return err
	}

	err = os.WriteFile(cacheDirPath+essencePath+fmt.Sprint(msg.Id), buffer.Bytes(), os.ModePerm)
	if err != nil {
		logger.Errorf("缓存群精华消息时出现错误: %v", err)
	} else {
		logger.Infof("缓存群精华消息成功: %d", msg.Id)
	}
	return err
}

func removeGroupEssence(msg int64) {
	err := os.Remove(cacheDirPath + essencePath + fmt.Sprint(msg))
	if err != nil {
		logger.Errorf("尝试移除群精华缓存消息 %d 时出现错误: %v", msg, err)
	}
}

func GetEssenceList(msg int64) []int64 {

	result := make([]int64, 0)

	files, err := os.ReadDir(cacheDirPath + essencePath)
	if err != nil {
		logger.Errorf("获取群精华缓存列表时错误: %v", err)
		return result
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		i, err := strconv.ParseInt(file.Name(), 0, 64)
		if err != nil {
			logger.Errorf("解析群精华缓存ID时出现错误: %v", err)
			continue
		}
		result = append(result, i)
	}

	return result
}

func FetchEssenceListToCache() (int, error) {
	gd, err := bot.Instance.GetGroupEssenceMsgList(ValGroupInfo.Code)
	if err != nil {
		return -1, err
	}
	logger.Infof("成功获取群精华消息: %d 则", len(gd))
	result := 0
	for _, digest := range gd {
		gpMsg, err := GetGroupMessage(digest.GroupCode, int64(digest.MessageID))
		if err != nil {
			logger.Errorf("尝试获取群精华消息 %d 时错误: %v", digest.MessageID, err)
			continue
		}
		if err = saveGroupEssenceErr(gpMsg); err == nil {
			result++
		}

	}
	return result, nil
}

// GetGroupEssenceMessage 获取瓦群群精华消息
func GetGroupEssenceMessage(msg int64) (result *message.GroupMessage, err error) {
	b, err := os.ReadFile(cacheDirPath + essencePath + fmt.Sprint(msg))

	if err == nil {
		persit := &PersistentGroupMessage{}
		buffer := bytes.NewBuffer(b)
		dec := gob.NewDecoder(buffer)
		err = dec.Decode(persit)
		if err != nil {
			logger.Errorf("群精华消息 %d 反序列化失败: %v", msg, err)
		} else {
			if result, err = persit.ToGroupMessage(); err == nil {
				fixGroupImages(ValGroupInfo.Code, result)
				logger.Infof("群精华消息 %d 获取成功.", msg)
			} else {
				logger.Errorf("群精华消息 %d 反序列化失败: %v", msg, err)
			}
		}
	}

	if result == nil {
		logger.Infof("尝试使用 QQ API 获取群精华消息 %d ...", msg)
		result, err = GetGroupMessage(ValGroupInfo.Code, msg)
		if err != nil {
			logger.Errorf("群精华消息 %d 获取失败: %v", msg, err)
		}
	}

	return
}
