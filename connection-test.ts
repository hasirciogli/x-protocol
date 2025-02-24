const REQUEST_BODY = {
  name: "hello",
  payload: {
    username: "jhon",
    password: "doe1234",
  },
};

// Performans ölçümü ve istatistikler için değişkenler
let successCount = 0;
let failureCount = 0;
let totalRequests = 0;
const TEST_DURATION = 1000 * 0.01;
const GROUP_LIMIT = 1; // Her grup için paralel istek sayısı

// POST istekleri gönderen fonksiyon
const sendRequests = async () => {
  const startTime = performance.now(); // Test başlama zamanı
  const endTime = startTime + TEST_DURATION; // Test bitiş zamanı
  let currentTime = startTime;
  const URL = "http://localhost:8080/calls"; // URL
  const failedRequests: { request: number; error: any }[] = [];

  // Test süresi boyunca istek göndermeye devam et
  while (currentTime < endTime) {
    const requests = Array.from({ length: GROUP_LIMIT }).map(
      async (_, index) => {
        totalRequests++; // Her istekte toplam isteği artır

        const requestStartTime = performance.now(); // İstek başlangıç zamanı
        try {
          const response = await fetch(URL, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              // Authorization: "Bearer 123456",
            },
            body: JSON.stringify(REQUEST_BODY),
          });

          // İstek başarılı ise
          if (response.ok) {
            successCount++; // Başarı sayısını artır
            console.log(
              `Request #${totalRequests + index} succeeded in ${(
                performance.now() - requestStartTime
              ).toFixed(2)} ms`
            );
          } else {
            console.log(await response.text());

            throw new Error(
              `Request failed with status ${response.status} - ${response.statusText} and body ${response.body}`
            );
          }
        } catch (error) {
          // İstek başarısızsa
          failureCount++; // Hata sayısını artır
          failedRequests.push({
            request: totalRequests + index,
            error: error as any,
          });
          console.log(
            `Request #${totalRequests} failed in ${(
              performance.now() - requestStartTime
            ).toFixed(2)} ms`
          );
        }
      }
    );

    // 5 istek gönderildikten sonra hepsinin tamamlanmasını bekleyin
    await Promise.all(requests);
    currentTime = performance.now(); // Geçen süreyi güncelle
  }

  // Test tamamlandıktan sonra istatistikleri yazdır
  const totalDuration = (performance.now() - startTime) / 1000; // Test süresi saniye cinsinden
  const successRate = (successCount / totalRequests) * 100;
  const failRate = (failureCount / totalRequests) * 100;

  console.log("Failed Requests:");
  console.log(failedRequests);

  console.log(`\nTest tamamlandı!`);
  console.log(`Test süresi: ${totalDuration.toFixed(2)} saniye`);
  console.log(`Toplam istek sayısı: ${totalRequests}`);
  console.log(`Başarı sayısı: ${successCount}`);
  console.log(`Hata sayısı: ${failureCount}`);
  console.log(`Başarı oranı: ${successRate.toFixed(2)}%`);
  console.log(`Hata oranı: ${failRate.toFixed(2)}%`);
};

// Testi başlat
sendRequests();
