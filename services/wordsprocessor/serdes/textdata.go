package serdes

import (
    "fmt"
    "google.golang.org/protobuf/proto"
    "example.com/textprocessor/services/wordsprocessor/generated/pb"
)

type TextDataSerde struct{}

func (s *TextDataSerde) SerializeObj(value interface{}) ([]byte, error) {
    v, ok := value.(*pb.TextData)
    if !ok {
        return nil, fmt.Errorf("value is not *pb.TextData")
    }
    return s.Serialize(v)
}

func (s *TextDataSerde) DeserializeObj(data []byte) (interface{}, error) {
    return s.Deserialize(data)
}

func (s *TextDataSerde) Serialize(value *pb.TextData) ([]byte, error) {
    return proto.Marshal(value)
}

func (s *TextDataSerde) Deserialize(data []byte) (*pb.TextData, error) {
    value := &pb.TextData{}
    if err := proto.Unmarshal(data, value); err != nil {
        return nil, err
    }
    return value, nil
}