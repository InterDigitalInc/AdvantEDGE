---
layout: default
title: EDGE Services
parent: Overview
nav_order: 3
has_children: true
permalink: docs/overview/edge-services/
---

Topic | Abstract
------|------
[Location Service](#location-service) | AdvantEDGE implementation of ETSI MEC013 specification that can be used by edge applications to gain location information on terminals and network access points.
[Radio Network Information Service](#radio-network-information-service) | AdvantEDGE implementation of ETSI MEC012 specification that can be used by edge applications to gain information about the cellular mobile network.
[Wireless Access Information Service](#wireless-access-information-service) | AdvantEDGE implementation of ETSI MEC028 specification that can be used by edge applications to gain information about the WLAN network.
[Edge Platform Application Enablement Service](#edge-platform-application-enablement-service) | AdvantEDGE implementation of ETSI MEC011 specification that can be used by edge applications to discover, advertise, consume and offer MEC services.
[Application Mobility Service](#application-mobility-service) | AdvantEDGE implementation of ETSI MEC021 specification that can be used by edge applications for relocation of user context and/or application instance across MEC platforms.
[Application State Transfer Service](#application-state-transfer-service) | AdvantEDGE proprietary implementation of a state transfer service that can be used by edge applications to transfer a state to another application instance.
NEXT STEP: [Platform APIs](#next-step) |

-----
## Location Service
AdvantEDGE provides a built-in Location Service implementation that integrates with scenarios.

This service provides the following capabilities:
- _Learning location of a device within the network_
- _Learning information on all devices located within a zone or connected to a point-of-access_
- _Getting real-time updates on device location as they move across the network_

Want to know more about Location service: [Location Service]({{site.baseurl}}{% link docs/overview/edge-services/overview-location-service.md %})

-----
## Radio Network Information Service
AdvantEDGE provides a built-in RNIS implementation that integrates with scenarios.

This service provides the following capabilities:
- _Learning radio network conditions_
- _Performing user plane measurements_
- _Learning about devices connected to the radio node(s) associated with the mobile edge host and their radio access bearers_
- _Getting real-time updates on devices radio conditions as they move across the network_

Want to know more about RNIS: [Radio Network Information Service]({{site.baseurl}}{% link docs/overview/edge-services/overview-rnis.md %})

-----
## Wireless Access Information Service
AdvantEDGE provides a built-in WAIS implementation that integrates with scenarios.

This service provides the following capabilities:
- _Learning terminal (station or STA) information_
- _Learning access-point information_
- _Getting real-time updates on devices WLAN conditions as they move across the network_

Want to know more about WAIS: [Wireless Access Information Service]({{site.baseurl}}{% link docs/overview/edge-services/overview-wais.md %})

-----
## Edge Platform Application Enablement Service
AdvantEDGE provides a built-in Edge Platform Application Enablement Service implementation that integrates with scenarios.

Mp1 reference point provides two different APIs: _MEC Application Support_ and _MEC Service Management_

These APIs allow MEC Applications to interact with the MEC System, such as:
- Application registration/deregistration
- Service discovery & offering
- Event notifications about service and application availability
- Traffic rules, DNS and time of day

Want to know more about App Enablement service: [Edge Platform Application Enablement Service]({{site.baseurl}}{% link docs/overview/edge-services/overview-app-enablement-service.md %})

-----
## Application Mobility Service
AdvantEDGE provides a built-in Application Mobility Service implementation that integrates with scenarios.

This service provides the following capabilities:
- TBD

Want to know more about AMS: [Application Mobility Service]({{site.baseurl}}{% link docs/overview/edge-services/overview-ams.md %})

-----
## Application State Transfer Service
AdvantEDGE provides a proprietary application state transfer service that facilitates state transfer between edge application instances.

This service provides the following capabilities:
- _Creating & configuring a Mobility Group (MG)_
- _Registering edge applications to the MG_
- _Executing application state transfers_

Want to know more about Application State Transfer: [Application State Transfer]({{site.baseurl}}{% link docs/overview/edge-services/overview-state-transfer.md %})

----
## Next Step
Learn about the various [Plarform APIs]({{site.baseurl}}{% link docs/overview/overview-api.md %}) that allows integration of your applications with the AdvantEDGE platform:
- Platform APIs
- Edge Service APIs
