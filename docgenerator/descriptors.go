package docgenerator

import (
	"strings"

	"github.com/nut-game/nano/constants"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func ProtoDescriptors(protoName string) ([]byte, error) {
	if strings.HasSuffix(protoName, ".proto") {
		// 处理 .proto 文件
		fileDescriptor, err := protoregistry.GlobalFiles.FindFileByPath(protoName)
		if err != nil {
			return nil, constants.ErrProtodescriptor
		}
		return fileDescriptor.(interface{ ProtoLegacyRawDesc() []byte }).ProtoLegacyRawDesc(), nil
	}

	// 处理消息类型
	if strings.HasPrefix(protoName, "types.") {
		protoName = strings.Replace(protoName, "types.", "google.protobuf.", 1)
	}

	// 查找消息类型
	messageType, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(protoName))
	if err != nil {
		return nil, constants.ErrProtodescriptor
	}

	// 获取消息的描述符
	descriptor := messageType.Descriptor()

	// 获取文件描述符
	fileDescriptor := descriptor.ParentFile()

	// 返回文件描述符的原始字节
	return fileDescriptor.(interface{ ProtoLegacyRawDesc() []byte }).ProtoLegacyRawDesc(), nil
}
