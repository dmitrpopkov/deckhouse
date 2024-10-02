---
title: "Сloud provider — OpenStack: настройки"
---

Модуль автоматически включается для всех облачных кластеров, развернутых в OpenStack.

{% include module-alerts.liquid %}

Количество и параметры процесса заказа машин в облаке настраиваются в custom resource [`NodeGroup`](../../modules/040-node-manager/cr.html#nodegroup) модуля `node-manager`, в котором также указывается название используемого для этой группы узлов инстанс-класса (параметр `cloudInstances.classReference` NodeGroup).  Инстанс-класс для cloud-провайдера OpenStack — это custom resource [`OpenStackInstanceClass`](cr.html#openstackinstanceclass), в котором указываются конкретные параметры самих машин.

Настройки модуля устанавливаются автоматически на основании выбранной схемы размещения. В большинстве случаев нет необходимости в ручной конфигурации модуля.

Если вам необходимо настроить модуль, потому что, например, у вас bare-metal-кластер, для которого нужно включить возможность добавлять дополнительные инстансы из OpenStack, смотрите раздел [Как настроить Hybrid-кластер в OpenStack](faq.html#как-поднять-гибридный-кластер).

> **Внимание!** При изменении настроек модуля **пересоздания существующих объектов `Machines` в кластере НЕ происходит** (новые объекты `Machine` будут создаваться с новыми параметрами). Пересоздание происходит только при изменении параметров `NodeGroup` и `OpenStackInstanceClass`. См. подробнее [в документации модуля node-manager](../../modules/040-node-manager/faq.html#как-пересоздать-эфемерные-машины-в-облаке-с-новой-конфигурацией).
> Для настройки аутентификации с помощью модуля `user-authn` необходимо в Crowd'е проекта создать новое `Generic` приложение.

## Список необходимых сервисов OpenStack

Список сервисов OpenStack, необходимых для работы Deckhouse Kubernetes Platform в OpenStack:

| Сервис                           | Версия API |
|:---------------------------------|:----------:|
| Identity (Keystone)              | v3         |
| Compute (Nova)                   | v2         |
| Network (Neutron)                | v2         |
| Block Storage (Cinder)           | v3         |
| Load Balancing (Octavia) &#8432; | v2         |

&#8432;  Если нужно заказывать Load Balancer.

В случае использования публичных облаков, адреса и порты API можно узнать в официальной документации сервиса:  
* [VK Cloud](https://cloud.vk.com/docs/tools-for-using-services/api/rest-api/endpoints)

{% include module-settings.liquid %}