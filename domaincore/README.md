# ddd core

Migrated from the **infr project**, extracting the common DDD components.

The current goal is to complete the infr project's refactoring, splitting it into four projects:
**project**, **platform**, **resource**, and **vault**,
while maintaining the same deployment model as Vercel.

The code will be revised and will not be exactly the same as the original **infr project** code.

## about `tags`

Added common `tags` to `Aggregate` and `BaseDomainEvent`. The domain is not pure, but it will save a lot of trouble.
