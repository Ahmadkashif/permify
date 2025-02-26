
<h1 align="center">
    <img src="https://raw.githubusercontent.com/Permify/permify/master/assets/permify-logo.svg" alt="Permify logo" width="336px" /><br />
    Permify - Open Source Authorization Service
</h1>

<p align="center">
    <a href="https://github.com/Permify/permify" target="_blank"><img src="https://img.shields.io/github/go-mod/go-version/Permify/permify?style=for-the-badge&logo=go" alt="Permify Go Version" /></a>&nbsp;
    <a href="https://goreportcard.com/report/github.com/Permify/permify" target="_blank"><img src="https://goreportcard.com/badge/github.com/Permify/permify?style=for-the-badge&logo=go" alt="Permify Go Report Card" /></a>&nbsp;
    <a href="https://github.com/Permify/permify" target="_blank"><img src="https://img.shields.io/github/license/Permify/permify?style=for-the-badge" alt="Permify Licence" /></a>&nbsp;
    <a href="https://discord.gg/MJbUjwskdH" target="_blank"><img src="https://img.shields.io/discord/950799928047833088?style=for-the-badge&logo=discord&label=DISCORD" alt="Permify Discord Channel" /></a>&nbsp;
    <a href="https://github.com/Permify/permify/pkgs/container/permify" target="_blank"><img src="https://img.shields.io/github/v/release/permify/permify?include_prereleases&style=for-the-badge" alt="Permify Release" /></a>&nbsp;
    <a href="https://img.shields.io/github/commit-activity/m/Permify/permify?style=for-the-badge" target="_blank"><img src="https://img.shields.io/github/commit-activity/m/Permify/permify?style=for-the-badge" alt="Permify Commit Activity" /></a>&nbsp;
</p>
<p align="center">
    <img src="https://raw.githubusercontent.com/Permify/permify/master/assets/permify-dsl.gif" alt="Permify - Open source authorization as a service" />
</p>

## What is Permify?
[Permify](https://github.com/Permify/permify) is an **open-source authorization service** for creating and maintaining fine-grained authorizations in your applications.

With Permify you can easily structure your authorization model, store authorization data in your own servers securely, and interact with Permify API to handle all authorization questions from any of your applications.

Permify's data model is inspired by Google’s consistent, global authorization system, [Google Zanzibar Paper](https://storage.googleapis.com/pub-tools-public-publication-data/pdf/41f08f03da59f5518802898f68730e247e23c331.pdf).

## Key Features

⚙️ **Production ready** authorization API that serve as **gRPC** and **REST**

🔮 Domain Specific Authorization Language - Permify Schema - to **easily model** your authorization

🔐 Database Configuration to store your permissions **in house** with **high availability**

✅ Ask authorization questions and get answers **down to 10ms** with **parallel graph engine**

💪 Battle tested, robust **authorization architecture and data model** based on [Google Zanzibar](https://storage.googleapis.com/pub-tools-public-publication-data/pdf/41f08f03da59f5518802898f68730e247e23c331.pdf)

📝 **Audit & Reason** your access control hassle-free with various functionalities via API

⚡ Analyze **performance and behavior** of your authorization with tracing tools [jaeger], [signoz] or [zipkin]

[jaeger]: https://www.jaegertracing.io/
[signoz]: https://signoz.io/
[zipkin]: https://zipkin.io/

### Permify works best:

- If you already have an identity/auth solution and want to plug in fine-grained authorization on top of that.
- If you want to create a unified access control mechanism to use across your individual applications.
- If you want to make future-proof authorization system and don't want to spend engineering effort for it.
- If you’re managing authorization for growing micro-service infrastructure.
- If your authorization logic is cluttering your code base.
- If your data model is getting too complicated to handle your authorization within the service.
- If your authorization is growing too complex to handle within code or API gateway.

## How it works

Permify stores access control relations on a database you choose and performs authorization checks based on the stored relations. We called that database Write Database (WriteDB) and it behaves as a centralized data source for your authorization system. You can model your authorization with Permify's DSL - Permify Schema and perform access checks with a single API call anywhere on your stack.

![Value Chain Schema](https://user-images.githubusercontent.com/34595361/186108668-4c6cb98c-e777-472b-bf05-d8760add82d2.png)

<p align="center">
    <a href="https://play.permify.co" target="_blank"><img src="https://raw.githubusercontent.com/Permify/permify/f548b4d5ae0d19c7d2f5bf61cc60ee04a58fa281/assets/Github%20Button.svg" alt="Try Permify" /></a>&nbsp;
</p>

## Getting Started

- [Set Up & Run Permify Service] in your environment.
- Follow a guide to model your authorization using [Permify Schema].
- Learn how Permify [Store Authorization Data] as relations.
- See our [playground], build your authorization logic and test it with sample data.
- Take a look at the overview of [Permify API].

[Set Up & Run Permify Service]: https://docs.permify.co/docs/installation/overview
[Store Authorization Data]: https://docs.permify.co/docs/getting-started/sync-data
[Permify Schema]: https://docs.permify.co/docs/getting-started/modeling
[playground]: https://play.permify.co/
[Permify API]: https://docs.permify.co/docs/api-overview/overview

[![Run in Postman](https://run.pstmn.io/button.svg)](https://god.gw.postman.com/run-collection/16122080-54b1e316-8105-4440-b5bf-f27a05a8b4de?action=collection%2Ffork&collection-url=entityId%3D16122080-54b1e316-8105-4440-b5bf-f27a05a8b4de%26entityType%3Dcollection%26workspaceId%3Dd3a8746c-fa57-49c0-83a5-6fcf25a7fc05)
[![View in Swagger](http://jessemillar.github.io/view-in-swagger-button/button.svg)](https://app.swaggerhub.com/apis-docs/permify/permify-api/v0.0.0-alpha4)

## Community & Support
We love to talk about authorization also we would love to hear from you :heart:

You can get immidiate help on our [Discord](https://discord.gg/MJbUjwskdH) channel. This can be any kind of questions related to Permify, authorization, or even from authentication or identity access control. We'd love to discuss anything related with access control space.

For feature requests, bugs or any improvements you can always open an issue. If you like Permify, please consider giving us a :star:

Want to contribute ? See [CONTRIBUTING.md](https://github.com/Permify/permify/blob/master/CONTRIBUTING.md).

<p align="left">
<a href="https://discord.gg/MJbUjwskdH">
 <img height="70px" width="70px" alt="permify | Discord" src="https://user-images.githubusercontent.com/39353278/187209316-3d01a799-c51b-4eaa-8f52-168047078a14.png" />
</a>
<a href="https://twitter.com/GetPermify">
  <img height="70px" width="70px" alt="permify | Twitter" src="https://user-images.githubusercontent.com/39353278/187209323-23f14261-d406-420d-80eb-1aa707a71043.png"/>
</a>
<a href="https://www.linkedin.com/company/permifyco">
  <img height="70px" width="70px" alt="permify | Linkedin" src="https://user-images.githubusercontent.com/39353278/187209321-03293a24-6f63-4321-b362-b0fc89fdd879.png" />
</a>
</p>

## Stargazers

[![Stargazers repo roster for permify/permify](https://reporoster.com/stars/permify/permify)](https://github.com/permify/permify/stargazers)
