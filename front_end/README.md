# معماری Social Downloader

## ساختار پوشه‌ها

```
lib/
├── data/                          # لایه Data (Model + Repository)
│   ├── models/
│   │   └── file_info.dart         # VideoInfo + FormatInfo (هم‌تراز با models.go)
│   └── repository/
│       ├── core/
│       │   └── api_client.dart    # Dio client
│       └── download_repository.dart
├── view/                          # لایه View
│   ├── shell/
│   │   └── main_shell.dart        # Bottom Navigation + StatefulNavigationShell
│   └── features/
│       ├── home/
│       │   ├── home_page.dart
│       │   └── download/          # BLoC این فیچر (نقش ViewModel)
│       │       ├── download_bloc.dart
│       │       ├── download_event.dart
│       │       └── download_state.dart
│       ├── history/
│       │   └── history_page.dart
│       └── settings/
│           └── settings_page.dart
├── router/
│   └── app_router.dart            # go_router + StatefulShellRoute
└── main.dart
```

## نقش‌ها در MVVM

| لایه MVVM      | در این پروژه                              | توضیح                               |
|---------------|-------------------------------------------|--------------------------------------|
| **Model**     | `data/models/` + Repository               | داده‌ها و ارتباط با بک‌اند           |
| **View**      | `view/features/*/..._page.dart` + Shell   | UI خالص، بدون منطق بیزنس           |
| **ViewModel** | `DownloadBloc` (و در آینده Blocهای دیگر) | مدیریت state و فراخوانی Repository |

BLoC در Flutter عملاً نقش ViewModel را بازی می‌کند: Event می‌گیرد، State تولید می‌کند و UI با `BlocBuilder` به آن واکنش نشان می‌دهد.

## Bottom Navigation

از `StatefulShellRoute.indexedStack` در go_router استفاده شده تا:

- state هر تب حفظ شود (با indexedStack)
- هر شاخه Navigator جدا داشته باشد
- مسیرها تمیز بمانند: `/` ، `/history` ، `/settings`

## صفحات جدید

1. **تاریخچه** (`HistoryPage`)  
   فعلاً آخرین نتیجه استخراج/دانلود موفق را از `DownloadBloc` نشان می‌دهد.  
   برای نسخه کامل: لیست تاریخچه را در local storage (مثلاً `shared_preferences` یا `hive`) یا endpoint بک‌اند نگه دارید.

2. **تنظیمات** (`SettingsPage`)  
   - کیفیت پیش‌فرض دانلود  
   - حالت تاریک (UI)  
   - آدرس بک‌اند  
   - اطلاعات نسخه  

   این‌ها فعلاً در state محلی صفحه هستند. برای پایداری می‌توانید یک `SettingsBloc` یا `SettingsCubit` + SharedPreferences اضافه کنید.

## مدل هم‌تراز با Go

`VideoInfo` و `FormatInfo` با `MediaInfo` و `Format` در `models.go` هم‌خوان شده‌اند (`media_type`, `formats`, ...).

## پیشنهادهای بعدی برای MVVM کامل‌تر

1. `SettingsCubit` / `HistoryCubit` جدا
2. ذخیره تاریخچه در local storage
3. تزریق وابستگی با `get_it` یا `riverpod` (اگر بخواهید از pure BLoC فراتر بروید)
4. انتقال `baseUrl` از تنظیمات به `ApiClient` به صورت پویا
