# Edge Application Deployment & Behavioral Models
Edge applications may be implemented following different *deployment* and *behavioral* models depending on the goal pursued.

In order to support different deployment & behavioral models, AdvantEDGE supports different types of [edge applications](./edge-applications.md) and an [application state transfer service](./use/app-state-transfer.md).


### Deployment Model
Edge application deployment model relates to the UE to edge application relationship

Below are a few deployment mopdel examples:

- Model 1: One edge application instance serving one UE (one-to-one relationship)<br></t>- multiple edge application instances can reside on the same physical node simultaneously

- Model 2: One edge application serving all UEs connected to a PoA (one to many relationship)<br></t>- one edge application instance resides on a physical node serving a localized geographical area

- Model 3: One edge application serving all UEs present in a Zone (one-to-many relationship)<br></t>- one edge application instance resides on a centrally located node serving a larger geographic area

- Other edge application deployment models may be valid too!

> *__What if your deployment model is not listed?__<br>
AdvantEDGE has been developped to provide as much flexibility as possible, so it may already support other deployment models not listed above out of the box.<br>
If your use case requires a deployment model that is not currently supported, AdvantEDGE can be extended to support it.*

### Behavioral Model
The behavioral model of an edge application can vary greatly depending on its  function, deployment model and overall design.

Below are some considerations that influence the behavioral model of an edge application:

- Bootstrapping<br></t>- When is the edge application instantiated?<br></t>- Does it follow an Always-available vs Just-in-time instantiation model?<br></t>- Where is the application bootstrapped?

- State Management<br></t>- Is the edge application stateful or stateless?<br><t>- Where do stateful applications get the inital UE state?<br></t>- Does the state need to be persisted when UE moves away?

- UE Mobility<br></t>- How does the edge application react to UE mobility events?<br><t>- Should edge application instance follow UE movement through the network?<br><t>- Should UE state be transferred to another instance when the UE moves through the network?<br></t>- Does the MEC platform provide an instance/state transfer service or does it happen at the application level (e.g. "over-the-top")

To help application developers & researchers with edge application design, AdvantEDGE allows to experiment with different models in an agile manner before any deployment happens on the real infrastructure.
