# Mikroservisler ve ara servisler için x-protocol

## 📖 Proje Hakkında

Bu proje, RPC akrabası olan bir çağrı sistemi geliştirmek amacıyla oluşturulmuştur. Kullanıcıların doğrudan RPC yazmalarına gerek kalmadan, Servisler arası 2 call tanımı ile kolayca iletişim kurabilirler. Proje, modüler yapısı ve kolay kullanımı ile dikkat çekmektedir. Ayrıca dışarıdan erişilemeyen sidecar ve ek ufak servisler için de proxy modeli de bulunmaktadır. Örnek dashboard dış servisinden -> cluster-controller -> sidecar şeklinde bir erişim tipinde controller-manager üzerinde hem kendi call-servisleri hem de sidecar için proxy katmanı eklenebilir bu durumda direkt olarak sidecar veya herhangi alt servisler için namespace ayırarak ona proxy katmanı ile erişebilirsiniz.

## 🚀 Başlangıç

### Gereksinimler

- Go 1.18 veya üzeri
- Node.js 16 veya üzeri
- pnpm

### GO Kurulum

1. **Go Modüllerini Yükleyin:**

   Modül yüklemeye gerek yoktur. Direkt olarak:

   ```go
    import "github.com/hasirciogli/x-protocol/packages/go/packages"
   ```

   şeklinde import edebilirsiniz.

2. **Örnek Server Oluşturun:**
   main.go dosyasını aşağıdaki gibi oluşturun.

   ```go
   package main

   import (
       "encoding/json"

       "github.com/hasirciogli/x-protocol/packages/go/packages"
   )

   type HelloPayload struct {
       Message string `json:"message"`
       Name    string `json:"name"`
   }

   func main() {
       server := packages.NewXProtocolServer("localhost", 8080)
       server.RegisterCall("hello", func(payload json.RawMessage) json.RawMessage {
           var p HelloPayload
           p.Message = "hello"
           p.Name = "world"

           str, err := json.Marshal(p)
           if err != nil {
               return json.RawMessage(`{"error": "` + err.Error() + `"}`)
           }
           return json.RawMessage(str)
       })
       server.Start()
   }
   ```

3. **Uygulamayı Başlatın:**
   main.go dosyasını aşağıdaki gibi çalıştırın.
   ```bash
   go run main.go
   ```

### NODEJS Kurulumu

1. **Nodejs Paketlerini Yükleyin:**

   ```bash
   pnpm add x-protocol
   ```

2. **Örnek Client Oluşturun:**

   ```typescript
   import { newXProtocolClient } from "x-protocol/dist/packages/client";

   const client = newXProtocolClient("localhost", 8080);
   const response = await client.call("hello", { name: "world" });
   console.log(response);
   ```

## 🛠️ Özellikler

| Özellik        | Açıklama                                 |
| -------------- | ---------------------------------------- |
| RPC Çağrıları  | RPC tabanlı çağrılar ile hızlı iletişim. |
| Modüler Yapı   | Her bir bileşen ayrı paketlerde.         |
| Farklı diller  | Go, TypeScript, Python, Java, C#, vb.    |
| Kolay Kullanım | Kullanıcı dostu arayüz ve dokümantasyon. |

## 📚 Kullanım

RPC çağrıları yapmak için aşağıdaki örnekleri inceleyebilirsiniz:

```go
// Go tarafında x-protocol çağrısı
response, err := client.Call("RegisteredCallAction", args)
if err != nil {
    log.Fatal(err)
}
```

```typescript
// TypeScript tarafında x-protocol çağrısı
const response = await client.call("RegisteredCallAction", args);
```

## 🎨 Katkıda Bulunma

Katkılarınızı bekliyoruz! Lütfen aşağıdaki adımları izleyerek katkıda bulunun:

1. Forklayın
2. Yeni bir dal oluşturun (`git checkout -b feature/YourFeature`)
3. Değişikliklerinizi yapın ve commit edin (`git commit -m 'Add some feature'`)
4. Dalınızı gönderin (`git push origin feature/YourFeature`)
5. Bir Pull Request açın

## 📄 Lisans

Bu proje [MIT Lisansı](LICENSE) altında lisanslanmıştır.

## 📞 İletişim

Herhangi bir sorunuz veya öneriniz varsa, lütfen benimle iletişime geçin:

- **E-posta:** mhasirciogli@gmail.com
- **GitHub:** [hasirciogli](https://github.com/hasirciogli)

---

**Teşekkürler! Projemi incelediğiniz için teşekkür ederiz!**

.

# Basit Bağlantı Stress Testi

Test yapılırken homelab kullanıldı cloudflare zerotrust ile dışarıya açılıp ağ gecikmesi ile beraber olası kopmalar göz önüne alınarak test yapıldı. Ayrıca tek makineden yapıldığı için hata oluştu, multi node ile test yapılsaydı veya proxy ile bu hatalar oluşmazdı.

connection-test.ts dosyasında test yapılmıştır. Dilerseniz kendi testlerinizi yapabilirsiniz.

```bash
npx tsx connection-test.ts
```

Test sonucu:

```bash
Test tamamlandı!
Test süresi: 241.42 saniye
Toplam istek sayısı: 17000
Başarı sayısı: 16996
Hata sayısı: 4
Başarı oranı: 99.98%
Hata oranı: 0.02%
```

## Farklı programlama dilleri arası test + proxy test :D

![x-protocol-multi-language-test](assets/images/x-protocol-multi-language-test.gif)
