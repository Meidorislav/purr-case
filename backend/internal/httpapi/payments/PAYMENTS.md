# Payments Notes For Inventory

## Что уже есть

- Есть `POST /payments/checkout`.
- Он защищён JWT auth.
- Фронт шлёт туда массив `items` в формате:

```json
{
  "items": [
    { "sku": "case_common", "quantity": 1 }
  ]
}
```

- Backend валидирует:
  - `items` не пустой
  - `sku` не пустой
  - `quantity > 0`
  - одинаковый `sku` два раза в одном запросе нельзя

- Дальше backend:
  - берёт `userId` из JWT
  - генерит локальный `orderId`
  - создаёт Xsolla payment token server-side
  - передаёт в Xsolla `purchase.items[]`
  - передаёт `settings.external_id = orderId`
  - возвращает клиенту `checkoutUrl`

То есть checkout уже рабочий.

## Что возвращает checkout

Пример ответа:

```json
{
  "orderId": "order-1234567890",
  "status": "new",
  "itemsCount": 1,
  "checkoutUrl": "https://sandbox-secure.xsolla.com/paystation4/?token=..."
}
```

## Что уже есть по webhook

- Есть `POST /payments/webhook`
- Xsolla webhook уже принимается
- подпись уже валидируется
- payload уже парсится

Сейчас из webhook уже можно достать:
- `notificationType`
- `status`
- `userId`
- `orderId`
- `transactionId`

Поддерживаются события:
- `user_validation`
- `payment`
- `order_paid`
- `order_canceled`
- `refund`

## Что важно для inventory

На стороне payments уже можно понять, что платёж успешный.

То есть после `order_paid` у нас уже есть данные, чтобы передать их в inventory-логику:
- кто купил (`userId`)
- какой заказ (`orderId`)
- какая транзакция (`transactionId`)

Но сама выдача предмета пользователю пока не сделана.

Inventory слой отвечает за:
- фактическую выдачу предмета пользователю
- защиту от повторной выдачи при повторном webhook
- при необходимости обработку `refund` / `order_canceled`

## Что нужно сделать дальше в inventory

Минимально:
- добавить метод, который может начислить предмет пользователю по `userId + sku + quantity`
- вызвать его после `order_paid`

В идеале:
- сделать идемпотентность
- не выдавать предмет дважды при повторном webhook
- хранить факт обработанной оплаты или транзакции

## ОГРАНИЧЕНИЯ СЕЙЧАС

- страна пользователя пока захардкожена
- отдельной серверной корзины нет
- payment / order state пока не сохраняется в отдельную таблицу

## По Inventory слою

Самый удобный вариант:
- inventory даёт метод сервиса вида “начислить предмет пользователю”
- payments вызывает его из webhook handler на `order_paid`