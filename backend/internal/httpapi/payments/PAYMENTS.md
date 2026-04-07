# Payments Notes For Inventory

## Что уже есть

- `POST /payments/checkout` защищён JWT auth.
- Фронт шлёт туда массив `items`:

```json
{
  "items": [
    { "sku": "case_common", "quantity": 1 }
  ]
}
```

- Backend валидирует `items`, генерит `orderId`, создаёт Xsolla payment token server-side и возвращает `checkoutUrl`.
- В Xsolla уходит `settings.external_id = orderId`.
- Этот же `orderId` сохраняется в локальной таблице `payment_orders`.
- Список купленных `sku + quantity` сохраняется в `payment_order_items`.

Пример ответа checkout:

```json
{
  "orderId": "order-1234567890",
  "status": "new",
  "itemsCount": 1,
  "checkoutUrl": "https://sandbox-secure.xsolla.com/paystation4/?token=..."
}
```

## Что уже есть по webhook

- `POST /payments/webhook` принимает webhook от Xsolla.
- Подпись Xsolla валидируется.
- Payload парсится.
- Поддерживаются события:
  - `user_validation`
  - `payment`
  - `order_paid`
  - `order_canceled`
  - `refund`

## Что теперь происходит на `order_paid`

Когда приходит `notification_type = "order_paid"`:

- payments берёт `external_id` из webhook
- находит локальный заказ в `payment_orders`
- достаёт его items из `payment_order_items`
- проверяет `processed_payment_events`, чтобы не обработать один webhook дважды
- вызывает inventory service
- inventory делает upsert в таблицу `inventory`

То есть предметы после оплаты теперь начисляются автоматически.

## Как работает защита от повторной выдачи

Для этого есть таблица `processed_payment_events`.

Если Xsolla повторно пришлёт тот же webhook:

- событие уже будет в `processed_payment_events`
- backend вернёт успешный ответ
- но предметы второй раз не начислит

## Что делает inventory слой

В inventory service добавлен метод начисления предметов:

- если `sku` у пользователя уже есть, увеличиваем `quantity`
- если `sku` ещё нет, создаём новую запись

Для этого в БД добавлен unique index на:

```sql
(user_id, sku)
```

## Что всё ещё не закрыто

- страна пользователя пока захардкожена
- отдельной серверной корзины нет
- `refund` / `order_canceled` пока не откатывают предметы
- нет отдельного API для просмотра payment/order status
