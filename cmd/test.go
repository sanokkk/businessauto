package main

import (
	"autoshop/internal/config"
	"autoshop/internal/storage"
	"autoshop/internal/storage/filters"
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

func main() {
	testMinioClient()
}

func testMinioClient() {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "BsdGKMYfWizG3KX5jXet"
	secretAccessKey := "TXVDFVbu4E10mfChKPFTTxprYd8JhB7Vbbdla2Im"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now setup

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	exist, err := minioClient.BucketExists(ctx, "default")
	if err != nil {
		log.Fatalln(err)
	}

	if !exist {
		log.Println("Дефолт бакета нет, создаю")
		if err := minioClient.MakeBucket(context.Background(), "default", minio.MakeBucketOptions{}); err != nil {
			log.Fatalln(err)
		}
	}

	//minioClient.GetObject(context.Background(), "default", "")

}

func testDb() {
	prodStore := storage.NewProductStore(&config.MustLoadConfig().DbConfig)

	filter := filters.ProductFilter{
		//TitleFilter: &filters.TitleFilter{Title: "first"},
		//MakerFilter: &filters.MakerFilter{Makers: []string{"Первый"}},
		PriceFilter:    &filters.PriceFilter{Min: 0, Max: 30},
		CategoryFilter: &filters.CategoryRangeFilter{Categories: []string{"Первая", "Вторая"}},
	}

	res, err := prodStore.GetWithFilter(filter, 0, 100, []filters.OrderBy{filters.OrderBy{Field: "Price", Desc: true}})

	if err != nil {
		panic(err)
	}

	fmt.Printf("res: %+v", res)
}
func upResume() {
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST", "hh.ru/applicant/resumes/touch", nil,
	)

	req.URL = &url.URL{Host: "hh.ru", Scheme: "https", Path: "/applicant/resumes/touch"}
	// добавляем заголовки
	req.Header.Add("cookie", "hhtoken=QXG65DHjdwLkjJ6GTsy42qvU4tCC; hhuid=isnk_y3Ojeu662X_n0NDjw--; _ym_uid=1711185731920602105; _ym_d=1711185731; hhul=338ed5f619c21f7f59d5512539523c4911a4e29dc9de7780743776776529d6a2; __ddg1_=hnSPDBmKUjdHuIyBqyuz; regions=1; region_clarified=NOT_SET; tmr_lvid=1a7600e4b1086d22ef3ed6633d4bc726; tmr_lvidTS=1711185731130; iap.uid=5901811324e54c7d8386751266e11bf1; _xsrf=25cbb6df923d8d189649325c14232abc; display=desktop; crypted_hhuid=229A903BDA4A921456918FA00309EE91E0666861B39F0C5FB522DE25E873A04A; GMT=3; _ym_isad=2; _ym_visorc=b; domain_sid=BzBsJchV5uedkV47JDObh%3A1724100236867; device_magritte_breakpoint=l; device_breakpoint=l; hhrole=applicant; _hi=146732243; crypted_id=867258D173F118A25A110384997A5583A5471C6F722104504734B32F0D03E9F8; __zzatgib-w-hh=MDA0dC0jViV+FmELHw4/aQsbSl1pCENQGC9LX3hecSImGEpbJHhYUQkmTBR7cylXPRIWcER2e108biRfOVURCxIXRF5cVWl1FRpLSiVueCplJS0xViR8SylEXE99JhkReXEjWA8NVy8NPjteLW8PKhMjZHYhP04hC00+KlwVNk0mbjN3RhsJHlksfEspNVZReSsdRXgmJk8LEV9GQ20xKm5uURtNF1BKWQh9KiEZem4oDQoPY3IzaWVpcC9gIBIlEU1HGEVkW0I2KBVLcU8cenZffSpCZyFgR1khSVVWCSYVe0M8YwxxFU11cjgzGxBhDyMOGFgJDA0yaFF7CT4VHThHKHIzd2UqP2skYklbI0ZHSWtlTlNCLGYbcRVNCA00PVpyIg9bOSVYCBI/CyYgE3tsI08JEVtGSW9vG382XRw5YxEOIRdGWF17TEA=+Rpvxw==; gsscgib-w-hh=g3sN0kJ20UAoGvUJy99oQqmZPSCbiRmgsmpSwpm0lIqXLOLw7zWa+ie/6C16iS3NRw6JU8YNpppXVh3cg5R/WCVQnbi7bhnwa/RW7fwM2jiY2ZZOs1RaNJJQQ5Rtr1s+g4wY67PapPwEsC6usBYe53JhyrgMM3d5FndbxiKx5bp6I/DL4HzSL+2EGq2pgGJMeBmvnmW2w8tYWCPD3ATHB9arytDr4w1H6sFZVpDVjk+KeTO5Y3xdZ9tuLONwnw==; cfidsgib-w-hh=FdiQd3YiEvwGzF7GnJXBfKrDs4RhJXo6yqMd/fqc+n6ZcAhGxvYnZL+v2VBNwihZZZwtH3ExbnXX/gVtWCOT9ZzC+S2lnxKxL+no3d00JNtN33qK9ANA7P9TLG+vJzz2MT7ukIllI7VGD1JGFnL10CkitZdkSFsuDavJPg==; cfidsgib-w-hh=FdiQd3YiEvwGzF7GnJXBfKrDs4RhJXo6yqMd/fqc+n6ZcAhGxvYnZL+v2VBNwihZZZwtH3ExbnXX/gVtWCOT9ZzC+S2lnxKxL+no3d00JNtN33qK9ANA7P9TLG+vJzz2MT7ukIllI7VGD1JGFnL10CkitZdkSFsuDavJPg==; gsscgib-w-hh=g3sN0kJ20UAoGvUJy99oQqmZPSCbiRmgsmpSwpm0lIqXLOLw7zWa+ie/6C16iS3NRw6JU8YNpppXVh3cg5R/WCVQnbi7bhnwa/RW7fwM2jiY2ZZOs1RaNJJQQ5Rtr1s+g4wY67PapPwEsC6usBYe53JhyrgMM3d5FndbxiKx5bp6I/DL4HzSL+2EGq2pgGJMeBmvnmW2w8tYWCPD3ATHB9arytDr4w1H6sFZVpDVjk+KeTO5Y3xdZ9tuLONwnw==; tmr_detect=0%7C1724100262955; fgsscgib-w-hh=ZLSp86f381ea94c16ac4e6195de7bf8ab63deb63") // добавляем заголовок Accept

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	//bytes, err := io.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
}

func login() {
	client := &http.Client{}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	writer.WriteField("accountType", "APPLICANT")
	writer.WriteField("username", "79999809197")
	writer.WriteField("password", "20020918aS")
	writer.WriteField("isBot", "false")

	err := writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}

	req, err := http.NewRequest(
		"POST", "hh.ru/applicant/resumes/touch", &buf,
	)

	req.URL = &url.URL{Host: "hh.ru", Scheme: "https", Path: "/account/login?backurl=/"}

	// добавляем заголовки
	req.Header.Add("cookie", "hhtoken=QXG65DHjdwLkjJ6GTsy42qvU4tCC; hhuid=isnk_y3Ojeu662X_n0NDjw--; _ym_uid=1711185731920602105; _ym_d=1711185731; hhul=338ed5f619c21f7f59d5512539523c4911a4e29dc9de7780743776776529d6a2; __ddg1_=hnSPDBmKUjdHuIyBqyuz; regions=1; region_clarified=NOT_SET; tmr_lvid=1a7600e4b1086d22ef3ed6633d4bc726; tmr_lvidTS=1711185731130; iap.uid=5901811324e54c7d8386751266e11bf1; _xsrf=25cbb6df923d8d189649325c14232abc; display=desktop; crypted_hhuid=229A903BDA4A921456918FA00309EE91E0666861B39F0C5FB522DE25E873A04A; GMT=3; _ym_isad=2; _ym_visorc=b; domain_sid=BzBsJchV5uedkV47JDObh%3A1724100236867; device_magritte_breakpoint=l; device_breakpoint=l; hhrole=applicant; _hi=146732243; crypted_id=867258D173F118A25A110384997A5583A5471C6F722104504734B32F0D03E9F8; __zzatgib-w-hh=MDA0dC0jViV+FmELHw4/aQsbSl1pCENQGC9LX3hecSImGEpbJHhYUQkmTBR7cylXPRIWcER2e108biRfOVURCxIXRF5cVWl1FRpLSiVueCplJS0xViR8SylEXE99JhkReXEjWA8NVy8NPjteLW8PKhMjZHYhP04hC00+KlwVNk0mbjN3RhsJHlksfEspNVZReSsdRXgmJk8LEV9GQ20xKm5uURtNF1BKWQh9KiEZem4oDQoPY3IzaWVpcC9gIBIlEU1HGEVkW0I2KBVLcU8cenZffSpCZyFgR1khSVVWCSYVe0M8YwxxFU11cjgzGxBhDyMOGFgJDA0yaFF7CT4VHThHKHIzd2UqP2skYklbI0ZHSWtlTlNCLGYbcRVNCA00PVpyIg9bOSVYCBI/CyYgE3tsI08JEVtGSW9vG382XRw5YxEOIRdGWF17TEA=+Rpvxw==; gsscgib-w-hh=g3sN0kJ20UAoGvUJy99oQqmZPSCbiRmgsmpSwpm0lIqXLOLw7zWa+ie/6C16iS3NRw6JU8YNpppXVh3cg5R/WCVQnbi7bhnwa/RW7fwM2jiY2ZZOs1RaNJJQQ5Rtr1s+g4wY67PapPwEsC6usBYe53JhyrgMM3d5FndbxiKx5bp6I/DL4HzSL+2EGq2pgGJMeBmvnmW2w8tYWCPD3ATHB9arytDr4w1H6sFZVpDVjk+KeTO5Y3xdZ9tuLONwnw==; cfidsgib-w-hh=FdiQd3YiEvwGzF7GnJXBfKrDs4RhJXo6yqMd/fqc+n6ZcAhGxvYnZL+v2VBNwihZZZwtH3ExbnXX/gVtWCOT9ZzC+S2lnxKxL+no3d00JNtN33qK9ANA7P9TLG+vJzz2MT7ukIllI7VGD1JGFnL10CkitZdkSFsuDavJPg==; cfidsgib-w-hh=FdiQd3YiEvwGzF7GnJXBfKrDs4RhJXo6yqMd/fqc+n6ZcAhGxvYnZL+v2VBNwihZZZwtH3ExbnXX/gVtWCOT9ZzC+S2lnxKxL+no3d00JNtN33qK9ANA7P9TLG+vJzz2MT7ukIllI7VGD1JGFnL10CkitZdkSFsuDavJPg==; gsscgib-w-hh=g3sN0kJ20UAoGvUJy99oQqmZPSCbiRmgsmpSwpm0lIqXLOLw7zWa+ie/6C16iS3NRw6JU8YNpppXVh3cg5R/WCVQnbi7bhnwa/RW7fwM2jiY2ZZOs1RaNJJQQ5Rtr1s+g4wY67PapPwEsC6usBYe53JhyrgMM3d5FndbxiKx5bp6I/DL4HzSL+2EGq2pgGJMeBmvnmW2w8tYWCPD3ATHB9arytDr4w1H6sFZVpDVjk+KeTO5Y3xdZ9tuLONwnw==; tmr_detect=0%7C1724100262955; fgsscgib-w-hh=ZLSp86f381ea94c16ac4e6195de7bf8ab63deb63") // добавляем заголовок Accept
	req.Header.Add("Content-Type", "multipart/form-data;")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println(resp.StatusCode)
	/*bytes, err := io.ReadAll(resp.Body)

	fmt.Println(string(bytes))*/
}

func testInterface(req interface{}) {
	categoryFilter := req.(filters.CategoryRangeFilter)

	fmt.Printf("%+v", categoryFilter)
}
