package bc

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ripemd160"
)

// Version 用于生成地址的版本 主网的版本
const Version = byte(0x00)

// TestVersion 用于生成测试网络地址的版本 m或者n开头
const TestVersion = byte(0x6F)

// P2SHVersion P2SH 类型的地址支持多重签名  3开头
const P2SHVersion = byte(0x05)

// AddressChecksumLen 用于生成地址的校验和位数
const AddressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() (*Wallet, error) {
	privateKey, publicKey, err := newKeyPair()
	if err != nil {
		return nil, err
	}
	return &Wallet{privateKey, publicKey}, nil
}

// newKeyPair 通过私钥创建公钥
func newKeyPair() (ecdsa.PrivateKey, []byte, error) {
	//1.椭圆曲线算法生成私钥
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return *privateKey, nil, err
	}
	////序列化私钥
	//marshalECPrivateKey, _ := x509.MarshalECPrivateKey(privateKey)
	//fmt.Println(string(Base58Encode(marshalECPrivateKey)))
	//2.通过私钥生成公钥
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey, nil
}

// GetAddress 获取钱包地址 根据公钥生成地址
func (wallet *Wallet) GetAddress() []byte {

	//1.使用RIPEMD160(SHA256(PubKey)) 哈希算法，取公钥并对其哈希两次
	ripemd160Hash := Ripemd160Hash(wallet.PublicKey)
	//2.拼接版本
	version_ripemd160Hash := append([]byte{Version}, ripemd160Hash...)
	//3.两次sha256生成校验和
	checkSumBytes := CheckSum(version_ripemd160Hash)
	//4.拼接校验和
	bytes := append(version_ripemd160Hash, checkSumBytes...)
	//5.base58编码
	return Base58Encode(bytes)
}

// GetTestAddress 获取钱包测试网地址 根据公钥生成地址
func (wallet *Wallet) GetTestAddress() []byte {

	//1.使用RIPEMD160(SHA256(PubKey)) 哈希算法，取公钥并对其哈希两次
	ripemd160Hash := Ripemd160Hash(wallet.PublicKey)
	//2.拼接版本
	version_ripemd160Hash := append([]byte{TestVersion}, ripemd160Hash...)
	//3.两次sha256生成校验和
	checkSumBytes := CheckSum(version_ripemd160Hash)
	//4.拼接校验和
	bytes := append(version_ripemd160Hash, checkSumBytes...)
	//5.base58编码
	return Base58Encode(bytes)
}
func (wallet *Wallet) GetPriKey() string {
	//序列化私钥
	marshalECPrivateKey, _ := x509.MarshalECPrivateKey(&wallet.PrivateKey)
	priBlock := pem.Block{
		Type:  "ECD PRIVATE KEY",
		Bytes: marshalECPrivateKey,
	}
	b := bytes.NewBuffer(make([]byte, 0))
	bw := bufio.NewWriter(b)
	err := pem.Encode(bw, &priBlock)
	if err != nil {
		panic(err)
	}
	bw.Flush()
	i := b.Bytes()
	return string(Base58Encode(i))
}
func (wallet *Wallet) ResetPriKey(key string) error {
	//反序列化私钥
	decode := Base58Decode([]byte(key))
	block, _ := pem.Decode(decode)
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	wallet.PrivateKey = *privateKey
	return nil
}

// Ripemd160Hash 将公钥进行两次哈希
func Ripemd160Hash(publicKey []byte) []byte {
	//1.hash256
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)

	//2.ripemd160
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)

	return ripemd160.Sum(nil)
}

// CheckSum 两次sha256哈希生成校验和
func CheckSum(bytes []byte) []byte {

	hash1 := sha256.Sum256(bytes)
	hash2 := sha256.Sum256(hash1[:])

	return hash2[:AddressChecksumLen]
}

// IsValidForAddress 判断地址是否有效
func (wallet *Wallet) IsValidForAddress(address []byte) bool {

	//1.base58解码地址得到版本，公钥哈希和校验位拼接的字节数组
	version_publicKey_checksumBytes := Base58Decode(address)
	//2.获取校验位和version_publicKeHash
	checkSumBytes := version_publicKey_checksumBytes[len(version_publicKey_checksumBytes)-AddressChecksumLen:]
	version_ripemd160 := version_publicKey_checksumBytes[:len(version_publicKey_checksumBytes)-AddressChecksumLen]

	//3.重新用解码后的version_ripemd160获得校验和
	checkSumBytesNew := CheckSum(version_ripemd160)

	//4.比较解码生成的校验和CheckSum重新计算的校验和
	if bytes.Compare(checkSumBytes, checkSumBytesNew) == 0 {
		return true
	}

	return false
}
