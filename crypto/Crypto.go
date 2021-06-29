package crypto

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
)

/**
des的CBC
编写填充函数,如果最后一个分组字节数不够,填充
如果最后一个分组字节数,也适合,相当于添加了一个新的分组
填充的每个字节的值=缺少的字节的数值
*/
func paddingLastGroup(plainText []byte, blockSize int) []byte {
	//求出最后一组中剩余的字节数
	padNum := blockSize - len(plainText)%blockSize

	//创建新的切片,长度是padNum 每个字节的值是byte(padNum)
	char := []byte{byte(padNum)}

	//切片创建,并初始化
	newPlain := bytes.Repeat(char, padNum)

	//新的数组追加到原始明文的后边
	newText := append(plainText, newPlain...)

	return newText
}

//拆除填充的数据
func unPaddingLastGroup(plainText []byte) []byte {
	//拿到切片最后一个字节
	length := len(plainText)
	lastChar := plainText[length-1]
	number := int(lastChar) //尾部填充的字节个数 //这里有个技巧当时在填充的时候填充的每个字节值记录了填充的字节数

	return plainText[:length-number]
}

//des cbc
func DesEncrypt(plainText, key []byte) []byte {
	//1.创建一个底层使用des的密码接口
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//2.使用的是cbc分组模式需要对明文分组进行填充
	newText := paddingLastGroup(plainText, block.BlockSize())

	//3.创建一个密码分组模式的接口对象
	iv := []byte("12345678") //初始化向量
	blockMode := cipher.NewCBCEncrypter(block, iv)

	//4.加密
	cipherText := make([]byte, len(newText))
	blockMode.CryptBlocks(cipherText, newText)

	return cipherText
}

//des cbc
func DesDecrypt(cipherText, key []byte) []byte {
	//1.创建一个底层使用des的密码接口
	block, err := des.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//2.创建一个密码分组模式的接口对象
	iv := []byte("12345678") //初始化向量
	blockMode := cipher.NewCBCDecrypter(block, iv)

	//解密
	plainText := make([]byte, len(cipherText))
	blockMode.CryptBlocks(plainText, cipherText)

	//拆除填充的数据
	return unPaddingLastGroup(plainText)
}

//aes ctr
func AesEncrypt(plainText, key []byte) []byte {
	//1.创建一个底层使用aes的密码接口
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//2.ctr不需要初始化向量 不需要对明文分组进行填充
	//创建一个密码分组模式的接口对象
	iv := []byte("12345678abcdefgh") //iv种子随机数 ase 16bytes
	stream := cipher.NewCTR(block, iv)

	//3.加密
	cipherText := make([]byte, len(plainText))
	stream.XORKeyStream(cipherText, plainText)

	return cipherText
}

//aes ctr
func AesDecrypt(cipherText, key []byte) []byte {
	//1.创建一个底层使用des的密码接口
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	//2. 创建一个密码分组模式的接口对象
	iv := []byte("12345678abcdefgh") //iv种子随机数 ase 16bytes
	stream := cipher.NewCTR(block, iv)

	//3.解密
	plainText := make([]byte, len(cipherText))
	stream.XORKeyStream(plainText, cipherText)

	return plainText
}

//对称加密
func TestDes() {
	key := []byte("abcd1234")
	src := []byte("如果明文刚好不需要填充,在尾部多填充一个分组,这个分组的每个字节值是16")

	cipherText := DesEncrypt(src, key)
	plainText := DesDecrypt(cipherText, key)
	fmt.Printf("%s \n", plainText)

	key2 := []byte("12345678abcdefgh")
	src2 := []byte("特点: 密文没有规律, 明文分组是和一个数据流进行的按位异或操作, 最终生成了密文")

	cipherText2 := AesEncrypt(src2, key2)
	plainText2 := AesDecrypt(cipherText2, key2)

	fmt.Printf("%s\n", plainText2)
}

//非对称加密
//RSA生成密钥对过程
func GenerateRsaKey(keySize int) {
	//RSA私钥生成过程
	//1.使用rsa中的GenerateKey方法生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		panic(err)
	}
	//2.通过x509标准将得到的rsa私钥序列化为ASN.1的DER编码字符串
	derText := x509.MarshalPKCS1PrivateKey(privateKey)
	//3.通过pem将设置好的数据进行编码,并写入磁盘文件
	//初始化一个pem.Block块
	block := pem.Block{
		Type:  "rsa private key", //这个地方写一个字符串就行
		Bytes: derText,
	}
	//out 准备一个文件指针
	outFile, err := os.Create("private.pem")
	if err != nil {
		panic(err)
	}
	//写入磁盘文件
	err = pem.Encode(outFile, &block)
	if err != nil {
		panic(err)
	}

	outFile.Close()

	//公钥生成过程
	//1.从得到的私钥对象中将公钥信息取出来
	publicKey := privateKey.PublicKey
	//2.通过x509标准得到rsa公钥序列化为字符串
	ders, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//3.将公钥字符串设置到pem格式块中
	block2 := pem.Block{
		Type:  "rsa public key", //这个地方写一个字符串就行
		Bytes: ders,
	}
	//4.通过pem将设置好的数据进行编码，并写入磁盘文件
	outFile, err = os.Create("public.pem")
	err = pem.Encode(outFile, &block2)
	if err != nil {
		panic(err)
	}
}

func RsaEncrypt(plainText []byte, fileName string) []byte {
	//1.将公钥文件中的公钥读出来,得到pem编码的字符串
	file, err := os.Open(fileName)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	buf := make([]byte, fileStat.Size())
	file.Read(buf)

	//2.pem解码
	block, _ := pem.Decode(buf)

	//3.使用x509将编码之后的公钥解析出来
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	//pubKey, bl := pub.(*rsa.PublicKey)//断言
	pubKey := pub.(*rsa.PublicKey)

	//4.使用得到的公钥通过rsa进行数据加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, plainText)
	if err != nil {
		panic(err)
	}

	return cipherText
}

func RsaDecrypt(cipherText []byte, fileName string) []byte {
	//1.将私钥文件中的私钥读出来,得到pem编码的字符串
	file, err := os.Open(fileName)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	buf := make([]byte, fileStat.Size())
	file.Read(buf)

	//2.pem解码
	block, _ := pem.Decode(buf)

	//3.使用x509将编码之后的私钥解析出来
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	//4.使用得到的私钥通过rsa进行数据解密
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	if err != nil {
		panic(err)
	}

	return plainText
}

func TestRsa() {
	//panic: crypto/rsa: message too long for RSA public key size
	//此时需要修改block的大小
	//GenerateRsaKey(1024)
	//GenerateRsaKey(4096)
	GenerateRsaKey(4096 * 2)
	src := []byte("百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫百度贴吧小爬虫")
	cipherText := RsaEncrypt(src, "public.pem")
	plainText := RsaDecrypt(cipherText, "private.pem")

	fmt.Printf("%s\n", plainText)
}

//单向散列函数又称为 消息摘要函数 哈希函杂凑函数

func myHash() {
	ha := sha256.New()
	src := []byte("单向散列函数又称为 消息摘要函数 哈希函杂凑函数")
	ha.Write(src)
	ha.Write(src)
	ha.Write(src)
	res := ha.Sum(nil)
	s := hex.EncodeToString(res)
	fmt.Printf("%s\n", s)
}

func myHash2() {
	res := sha256.Sum256([]byte("单向散列函数又称为 消息摘要函数 哈希函杂凑函数"))
	fmt.Printf("%x\n", res)
}

//生成消息认证码
func GenerateHmac(plainText, key []byte) []byte {
	//1.创建哈希接口,需要指定使用哈希算法和秘钥
	hi := hmac.New(sha1.New, key)
	//2. 给哈希对象添加数据
	hi.Write(plainText)
	//3.计算散列值
	hashText := hi.Sum(nil)
	return hashText
}

//校验消息认证码
func VerifyHmac(plainText, key, hashText []byte) bool {
	//1.创建哈希接口,需要指定使用哈希算法和秘钥
	hi := hmac.New(sha1.New, key)
	//2. 给哈希对象添加数据
	hi.Write(plainText)
	//3.计算散列值
	newHashText := hi.Sum(nil)
	//4.通过新的散列值 和 接收的散列值进行比较
	return hmac.Equal(newHashText, hashText)
}

func TestHmac() {
	src := []byte("1.创建哈希接口,需要指定使用哈希算法和秘钥 2. 给哈希对象添加数据 3.计算散列值 4.通过新的散列值 和 接收的散列值进行比较")
	key := []byte("your key")

	hamc1 := GenerateHmac(src, key)
	bl := VerifyHmac(src, key, hamc1)
	log.Printf("%t", bl)
}

//RSA进行数字签名 使用私钥进行数字签名
func SignatureRSA(plainText []byte, privateKeyFileName string) []byte {
	//打开磁盘的私钥文件
	file, err := os.Open(privateKeyFileName)
	if err != nil {
		panic(err)
	}
	//将私钥文件的内容读出来
	info, err := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	file.Close()
	//使用pem对数据解码,得到pem.Block结构体变量
	block, _ := pem.Decode(buf)
	//x509将数据解析成私钥结构体,得到了私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//创建一个哈希对象 -> md5/sha1/sha512
	hi := sha512.New()
	//给对象添加数据
	hi.Write(plainText)
	//计算哈希值
	hashText := hi.Sum(nil)
	//使用rsa中的函数对散列值签名
	sigText, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA512, hashText)
	if err != nil {
		panic(err)
	}
	return sigText
}

//RSA验证数字签名
//使用公钥进行数字认证
func SignatureVerifyRSA(plainText []byte, publicKeyFileName string, sigText []byte) bool {
	//打开磁盘的公钥文件
	file, err := os.Open(publicKeyFileName)
	if err != nil {
		panic(err)
	}
	//将公钥文件的内容读出来 []byte
	info, err := file.Stat()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, info.Size())
	file.Read(buf)
	file.Close()
	//使用pem解码-得到pem.Block结构体变量
	block, _ := pem.Decode(buf)
	//使用x509对pem.Block中的Bytes变量中的数据进行解析 -> 得到一个接口
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(fmt.Sprintf("err 1 : %s", err))
	}
	//进行类型断言 ->得到了公钥结构体
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对原始消息进行哈希运算(和签名使用的哈希算法一致) -> 散列值
	hashText := sha512.Sum512(plainText)
	//签名认证 - rsa中的函数
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, hashText[:], sigText)
	if err == nil {
		return true
	}

	return false
}

func TestRsaSig() {
	GenerateRsaKey(1024)
	src := []byte("wo zai ping duo duo shang mai le hen duo de hao chi de ")
	sigText := SignatureRSA(src, "./private.pem")
	log.Println(sigText)
	bl := SignatureVerifyRSA(src, "./public.pem", sigText)
	log.Printf("%t", bl)
}

//椭圆曲线秘钥对的生成
func GenerateEccKey() {
	//1.使用ecdsa来生成秘钥对
	privateKey, err := ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	if err != nil {
		panic(err)
	}

	//2.将私钥写入磁盘
	//- 使用x509进行序列化
	derText, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	//- 将得到的切片字符串放入pem.Block结构体中
	block := pem.Block{
		Type:  "ecdsa private key",
		Bytes: derText,
	}
	//-使用pem编码
	file, err := os.Create("ecc_private.pem")
	if err != nil {
		panic(err)
	}

	err = pem.Encode(file, &block)
	if err != nil {
		panic(err)
	}
	file.Close()

	//3.将公钥写入磁盘
	//-从私钥中得到公钥
	publicKey := privateKey.PublicKey
	//使用x509进行序列化
	derText, err = x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//将得到的切片字符串放入pem.Block结构体中
	block = pem.Block{
		Type:  "ecdsa public key",
		Bytes: derText,
	}
	//使用pem编码
	file, err = os.Create("ecc_public.pem")
	if err != nil {
		panic(err)
	}
	pem.Encode(file, &block)
	file.Close()
}

func EccSignature(plainText []byte, privateKeyFileName string) (rText, sText []byte) {
	//打开私钥文件,将内容读出来 -》 []byte
	file, err := os.Open(privateKeyFileName)
	if err != nil {
		panic(err)
	}
	info, err := file.Stat()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, info.Size())
	file.Read(buf)
	file.Close()
	//使用pem进行数据解码 -> pem.Decode()
	block, _ := pem.Decode(buf)
	//使用x509,对私钥进行还原
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//对原始数据进行哈希运算 -》 散列值
	hashText := sha1.Sum(plainText)
	//对散列值进行数字签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashText[:])
	if err != nil {
		panic(err)
	}
	//得到的r和s不能直接使用,因为这是指针
	//应该将这两块内存中的数据进行序列化 -》 []byte
	rText, err = r.MarshalText()
	if err != nil {
		panic(err)
	}
	sText, err = s.MarshalText()
	if err != nil {
		panic(err)
	}

	return
}

func EccVerifySignature(plainText []byte, publicKeyFileName string, rText, sText []byte) bool {
	//打开公钥文件,将内容读出来 -》 []byte
	file, err := os.Open(publicKeyFileName)
	if err != nil {
		panic(err)
	}
	info, err := file.Stat()
	if err != nil {
		panic(err)
	}
	buf := make([]byte, info.Size())
	file.Read(buf)
	file.Close()
	//使用pem进行数据解码 -> pem.Decode()
	block, _ := pem.Decode(buf)
	//使用x509,对公钥进行还原
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//将接口转换成公钥
	//类型转换 断言
	publicKey := publicKeyInterface.(*ecdsa.PublicKey)
	//对原始数据进行哈希运算 -》 散列值
	hashText := sha1.Sum(plainText)
	//签名的验证
	var r, s big.Int
	r.UnmarshalText(rText)
	s.UnmarshalText(sText)
	return ecdsa.Verify(publicKey, hashText[:], &r, &s)
}

func TestEccSig() {
	GenerateEccKey()
	src := []byte("ECC sig ")
	rText, sText := EccSignature(src, "ecc_private.pem")
	b1 := EccVerifySignature(src, "ecc_public.pem", rText, sText)
	log.Printf("%t", b1)
}

func RunEcdsa() {
	cure := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(cure, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pubKey := privateKey.PublicKey

	data := "hello world!"
	hash := sha256.Sum256([]byte(data))

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("pubkey: %v\n", pubKey)
	fmt.Printf("r: %v len = %d\n", r.Bytes(), len(r.Bytes()))
	fmt.Printf("s: %v len = %d\n", s.Bytes(), len(s.Bytes()))

	//把 r, s进行序列化传输
	signature := append(r.Bytes(), s.Bytes()...)

	//....传输

	//定义两个辅助的big.Int
	r1 := big.Int{}
	s1 := big.Int{}
	//拆分我们signature,平均分,前半部分r,后半部分s
	m := len(signature) / 2
	//r1.SetBytes(signature[0:m])
	r1.SetBytes(signature[:m])
	s1.SetBytes(signature[m:])

	//校验需要三个东西: 数据, 签名, 公钥
	res := ecdsa.Verify(&pubKey, hash[:], &r1, &s1)
	fmt.Printf("校验结果 %t\n", res)

}
