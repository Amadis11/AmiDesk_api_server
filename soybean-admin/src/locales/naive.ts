import { dateEnGB, dateEnUS, dateItIT, dateRuRU, enGB, enUS, itIT, ruRU } from 'naive-ui';
import type { NDateLocale, NLocale } from 'naive-ui';

export const naiveLocales: Record<App.I18n.LangType, NLocale> = {
  'en-US': enUS,
  'en-GB': enGB,
  'ru-RU': ruRU,
  'it-IT': itIT
};

export const naiveDateLocales: Record<App.I18n.LangType, NDateLocale> = {
  'en-US': dateEnUS,
  'en-GB': dateEnGB,
  'ru-RU': dateRuRU,
  'it-IT': dateItIT
};
