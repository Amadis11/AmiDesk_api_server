# AmiDesk API Server - Amitronic 2.0

## Cel i zakres

Wydanie Amitronic 2.0 przygotowuje AmiDesk API Server do bezpiecznego wdrożenia ze stockowym klientem RustDesk. Zakres obejmuje automatyczne przekazywanie publicznego klucza hbbs przez HTTPS, wzmocnienie logowania do panelu administracyjnego i kont API, angielskie ustawienia domyślne oraz powtarzalny obraz Docker.

Wydanie nie modyfikuje protokołu zwykłych połączeń RustDesk opartych na identyfikatorze urządzenia i haśle.

## Architektura

Rozwiązanie składa się z trzech niezależnych elementów:

1. `hbbs` i `hbbr` obsługują rejestrację, zestawianie połączeń i relay RustDesk.
2. AmiDesk API Server obsługuje konta, panel administracyjny, urządzenia, adres book oraz endpoint heartbeat.
3. Reverse proxy udostępnia panel i API wyłącznie przez HTTPS.

AmiDesk API Server odczytuje wyłącznie publiczny plik `id_ed25519.pub`. Prywatny klucz `id_ed25519` pozostaje wyłącznie po stronie hbbs i nie może być montowany w kontenerze API, zwracany przez endpointy ani zapisywany w logach.

## Bootstrap stockowego klienta

Stockowy klient RustDesk wysyła anonimowy `POST /api/heartbeat`. Klient nie przesyła w tym żądaniu tokenu konta, dlatego bootstrap ma charakter anonimowy i jest chroniony transportowo przez HTTPS.

Po otrzymaniu heartbeat serwer przekazuje następujące opcje:

- `key` - publiczny klucz hbbs;
- `custom-rendezvous-server` - serwer ID;
- `relay-server` - serwer relay;
- `api-server` - adres API.

Klient zachowuje lokalnie wersję strategii. Serwer zwraca konfigurację tylko po zmianie konfiguracji serwerów lub zamontowanego publicznego klucza. Dzięki temu obrót klucza jest propagowany bez niepotrzebnego zapisu ustawień przy każdym heartbeat.

### Klient z samym kluczem hbbs

Klient mający ustawione ID Server, Relay Server i publiczny klucz, lecz bez API Server, nadal korzysta z podstawowej ścieżki RustDesk `ID + hasło`. Samo wdrożenie API nie zmienia jego lokalnej konfiguracji ani nie odcina go od hbbs/hbbr.

Taki klient nie może jednak korzystać z logowania konta, synchronizacji książki adresowej ani bootstrapu API, dopóki administrator lub użytkownik nie ustawi adresu API. Po jednorazowym ustawieniu API Server klient odbierze przez heartbeat pełną konfigurację i zapisze ją lokalnie.

### Klient zalogowany do API

Logowanie do API pozostaje opcjonalne dla zwykłych połączeń peer-to-peer. Po zalogowaniu klient otrzymuje token konta, wykorzystywany przez funkcje konta oraz przez część ścieżek TCP klienta. To nie zmienia ważności wcześniej ustawionego klucza ani adresów hbbs/hbbr.

## Secure TCP i decyzja o forku hbbs

Samo przekazanie klucza publicznego nie wymusza secure TCP dla anonimowego ruchu relay. W stockowym kliencie secure TCP jest inicjowane przy niepustym kluczu oraz tokenie konta albo `switch_code`.

Wdrożenie API nie wymaga forka hbbs, aby obsługiwać zwykłe, niezalogowane połączenia. Nie gwarantuje jednak połączeń mieszanych: klient zalogowany do API ma token i może inicjować secure TCP, podczas gdy drugi klient bez API pozostaje w ścieżce bez tokenu. Jeżeli hbbs nie obsługuje zgodnego z klientem KeyExchange dla tej ścieżki tokenowej, połączenie przez relay między tymi klientami może się nie zestawić.

Do czasu pozytywnego testu połączeń mieszanych albo wdrożenia zgodnego forka hbbs nie należy uznawać scenariusza „klient zalogowany do API -> klient bez API” za wspierany. Bezpośrednie połączenie P2P może zadziałać zależnie od NAT, ale nie jest to zastępstwo dla poprawnej obsługi relay.

Przed wdrożeniem forka należy wykonać testy z trzema przypadkami:

1. klient bez konta i z publicznym kluczem;
2. klient zalogowany do API, z tokenem i publicznym kluczem;
3. urządzenie docelowe bez konta API.

Fork hbbs musi zachować obsługę tokenless `PunchHoleRequest` oraz implementować framing KeyExchange oczekiwany przez stockowego klienta. API nie powinno tworzyć własnego tokenu między API i hbbs ani zmieniać protokołu rendezvous.

## Zmiany bezpieczeństwa

- Publiczny klucz hbbs jest dystrybuowany anonimowo wyłącznie przez HTTPS.
- Endpoint administracyjny pokazuje status bootstrapu i fingerprint klucza, ale nie jego treść.
- Sekrety TOTP nie są zwracane przez listę użytkowników panelu administracyjnego.
- Dla administratorów dostępne jest wymuszenie TOTP przez `requireAdminTOTP`. Ustawienie pozostaje wyłączone do czasu skonfigurowania TOTP dla każdego aktywnego administratora.
- Logowania do API i panelu są blokowane po pięciu nieudanych próbach dla tego samego konta na 15 minut. Licznik nie zależy od adresu IP, więc zmiana adresu przez automaty nie omija blokady.
- Udane logowanie czyści licznik błędów. Odpowiedzi nie rozróżniają nieistniejącego użytkownika i nieprawidłowego hasła.

Limit prób jest utrzymywany w pamięci procesu. Dla wdrożenia z wieloma replikami wymagany jest wspólny magazyn stanu, na przykład Redis.

## Ustawienia domyślne

- Strefa czasowa bazy danych i kontenera: `Europe/Warsaw`.
- Domyślny język panelu: English (UK).
- Chińska lokalizacja, chińskie komunikaty aktywnego panelu i zawartość sponsoringu upstream zostały usunięte.
- Widok Home przedstawia changelog Amitronic `v2.0` oraz angielski wpis historyczny kompatybilności klienta.

## Docker i trwałość konfiguracji

Obraz Docker buduje backend Go oraz panel Vue w osobnych etapach. Instalacja pnpm używa lockfile w trybie zamrożonym, z jawnym zatwierdzeniem wyłącznie wymaganego skryptu `esbuild`.

Przy pierwszym uruchomieniu skrypt startowy kopiuje przygotowaną konfigurację do trwałego wolumenu danych. Kolejne uruchomienia nie nadpisują tej konfiguracji. Skrypt ma zakończenia linii LF, wymagane przez `/bin/sh` w Alpine.

W obrazie kontenera katalog konfiguracji jest ustawiony przez `RUSTDESK_API_CONFIG_DIR=/app/data`. Polecenia CLI uruchomione przez `docker exec`, w tym dodawanie użytkownika, korzystają z tej samej trwałej konfiguracji i bazy co działająca usługa, niezależnie od domyślnego katalogu roboczego kontenera.

Wdrożenie musi:

1. przechowywać trwały wolumen danych;
2. montować `id_ed25519.pub` jako plik tylko do odczytu;
3. używać własnej wartości `signKey`;
4. udostępniać API przez reverse proxy HTTPS;
5. nie wystawiać bezpośrednio portu aplikacji do Internetu.

## Weryfikacja wykonana dla wersji 2.0

- Pełne testy backendu Go zakończone powodzeniem: `go test ./...`.
- Produkcyjny build panelu zakończony powodzeniem: `pnpm build`.
- Lokalny obraz Docker został zbudowany poprawnie.
- Kontener uruchomił API, utworzył trwałą konfigurację i odpowiedział `200` na heartbeat.
- Test bootstrapu z poprawnym przykładowym kluczem publicznym potwierdził zwrot klucza oraz adresów ID, relay i API w `config_options`.

## Kolejność wdrożenia produkcyjnego

1. Udostępnić API przez HTTPS i skonfigurować właściwe adresy ID, relay i API.
2. Zamontować prawdziwy publiczny klucz hbbs oraz potwierdzić fingerprint w panelu administratora.
3. Przetestować świeżego stockowego klienta bez konta, następnie klienta zalogowanego i urządzenie bez konta API.
4. Skonfigurować oraz zweryfikować TOTP dla każdego administratora.
5. Włączyć `requireAdminTOTP`.
6. Potwierdzić przez relay połączenie klienta zalogowanego do API z klientem bez API. Jeżeli test nie przejdzie, wdrożyć zgodny fork hbbs przed udostępnieniem logowania API użytkownikom.