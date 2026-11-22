# Competitive Advantage — Mimamori (2026-11-22)

Method: competitor survey (paid and free), money-flow decomposition, then
the Musk method + Socratic questioning repeated until a premise broke. Two
full rounds were run; each round records what was questioned, what was
deleted, and what survived. Judged from the maker's principles: the maker
holds no personal data, provides software only (no operating cost), and
the product itself is free.

## The market, and where the money goes

| Product | Form | Price signal |
|---|---|---|
| Fing | App + cloud device-identification DB | Free: 3 scans/day, ads. Premium $79.99/yr, Professional $179.99/yr. Monitoring needs Fing Desktop/Agent resident, or a Fingbox |
| NETGEAR Armor (Bitdefender) | Router suite | $99.99/yr (Plus $149.99/yr), on top of the router |
| eero Plus (Amazon) | Router suite | $99.99/yr |
| TP-Link HomeShield | Router suite | Basic free; Security+ $35.99/yr, parental $17.99/yr, Total Security $129.99/yr renewal — stacked add-ons |
| ASUS AiProtection (Trend Micro) | Router suite | "Free for life", recovered in the router's price and the vendor partnership |
| Avast Network Inspector | Bundled in AV suite | Scan is free; continuous monitoring requires a paid Premium subscription |
| Firewalla | Dedicated appliance | $189–949 one-time, no subscription — but an always-on box plus app/cloud pairing |
| Nmap / Angry IP Scanner | Free, open source | Lists and ports only; no verdict, no fixes |

Decomposing the money, a household pays for five things, and diagnosis is
none of them:

1. **Monitoring rent** — Fing, Armor, eero Plus. What is sold is not a
   check but a continuing relationship: alerts, timelines, retention.
2. **Router lock-in** — the suites run only on one maker's routers, so
   the security layer becomes a switching cost at replacement time.
3. **Always-on hardware margin** — Firewalla, Fingbox.
4. **Suite upsell** — the free scan is the funnel; the subscription
   suite is the product (Avast's monitoring paywall is the clearest
   form).
5. **Data** — Fing's identification database is built from user scans;
   AV telemetry follows the same shape. The scanner is itself a
   collector on your LAN.

Every competitor's charging point is a *continuing relationship*: an
account, a cloud, a resident agent, dedicated hardware, or a router
brand.

## The iteration

### Round 1 — question "always-on" and "more checks"

- Questioned premise: home safety needs continuous monitoring. Deleted:
  the risks Mimamori checks are configuration states (an open Telnet, a
  passwordless MQTT). They change only when a setting changes or a
  device is added, so a single run catches them at the moment of
  asking. New-device notice also survives deletion of the resident
  agent: it falls out of the history diff on each run. "Always-on" is
  the vendors' requirement — it is what makes the rent meterable — not
  the user's.
- Questioned premise: more checks is better. Deleted: comprehensiveness
  (expert tools already own it). Survived: verifying actual
  unauthenticated access with read-only handshakes, and presenting a
  verdict — a score, a plain-language fix, the points recovered.
- "Free is the advantage" — half deleted: router-bundled free tiers and
  Avast's free scan exist. What survived is free *without* hardware
  lock-in, an account, or an upsell funnel behind it.

### Round 2 — question "the product is visibility" (the premise broke)

- The industry's core product is visibility: device inventories,
  dashboards, timelines. Visibility is what can be metered — per
  network, per month, per retention window — so visibility is what gets
  sold.
- But the household's job is a verdict: "am I safe right now, and what
  do I fix?" A verdict cannot be metered; the moment it is delivered,
  the transaction is over. No subscription business can make the
  verdict its core product, because its model needs the relationship
  to continue.
- Second break: "a security tool needs a cloud and an account."
  Mimamori's checks are configuration checks and need neither — and
  for a tool whose job is to inspect your home, being offline and
  auditable from source is itself the trust mechanism. Incumbents sit
  on the other side: their scanner collects the home's device data,
  and giving that up would break their data asset and their funnel.

## The structural advantage

Mimamori can keep handing over home-network safety as a one-time
verdict instead of a continuing, metered relationship. Every
competitor's revenue rests on the relationship itself — subscription,
account, always-on hardware, router lock-in — so a form that delivers
the verdict and the fixes unconditionally, with no account, no cloud
and nothing left running, cannot be copied without breaking the model
that pays them. Mimamori has no revenue and no server, so it can keep
that promise permanently.

Supporting facts already true of the product:

- It never contacts the Internet and sends no telemetry; the only
  persistent data is a device-history file on the user's own PC.
- Standard library only, a single static binary: the "collects
  nothing" claim is auditable from source, not a policy page.
- It verifies unauthenticated access with read-only probes instead of
  listing open ports, and changes nothing on the network it checks.

## Remaining opportunities (ranked)

1. State the difference at the top of the README (no account, no
   subscription, nothing left running; a verdict, not a dashboard).
   Done alongside this document.
2. Make "check another home you look after" a first-class story: a
   single binary on a USB stick, no install, no trace left, a separate
   `--state` file per home. The README already documents the mechanics;
   a short guide in docs/ would turn them into the distribution wedge
   no subscription product can follow.
3. Make the offline claim checkable by the user: reproducible builds
   and a documented way to watch Mimamori make no external connections
   turn "trust us" into "verify it".
4. Collect the scheduled-run recipes (Task Scheduler / cron, `--json`,
   exit codes, the `incomplete` flag) into one docs/ guide, so a
   monthly checkup habit needs no resident software at all.
