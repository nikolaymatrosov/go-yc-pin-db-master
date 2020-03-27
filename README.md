### Принцип работы

### Deploy
#### MacOs (Linux)
Предполагается что у вас уже настроены [yc](https://cloud.yandex.ru/docs/cli/quickstart) и [s3cmd](https://cloud.yandex.ru/docs/storage/instruments/s3cmd). Они понадобятся для скрипта деплоя.

Чтобы задеплоить функции в ваше облако выполните:
1. Создать файл `.env` из шаблона из `.env.template` и заполнить
1. `./script/create.sh` создаст все необходимые объекты в облаке: 3 функции, очередь сообщений и 3 триггера.
1. `./script/deploy.sh`

#### Windows

Deploy script можно запустить в mingw (git bash).
Для его корректной работы вам может понадобиться
1. Скачать и установить [GnuZip](http://gnuwin32.sourceforge.net/packages/zip.htm)
1. Прописать `GnuZip` в path.

Возможно потребуются права администаратора.
