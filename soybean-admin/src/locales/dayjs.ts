import { locale } from 'dayjs';
import 'dayjs/locale/en';
import 'dayjs/locale/en-gb';
import 'dayjs/locale/ru';
import 'dayjs/locale/it';
import { localStg } from '@/utils/storage';

/**
 * Set dayjs locale
 *
 * @param lang
 */
export function setDayjsLocale(lang: App.I18n.LangType = 'en-GB') {
  const localMap = {
    'en-US': 'en',
    'en-GB': 'en-gb',
    'ru-RU': 'ru',
    'it-IT': 'it'
  } satisfies Record<App.I18n.LangType, string>;

  const l = lang || localStg.get('lang') || 'en-GB';

  locale(localMap[l]);
}
