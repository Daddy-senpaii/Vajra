# Git workflow

1. Keep `main` as the integration branch.
2. Start each feature from the latest `main` using a dedicated branch such as `feature/authentication` or `feature/running-mode`. Infrastructure changes use `infrastructure/<topic>`; fixes use `bugfix/<topic>`.
3. Make focused commits using Conventional Commit subjects (`feat:`, `fix:`, `docs:`, `ci:`, etc.).
4. Push the branch and open a Pull Request. Do not push feature implementation directly to `main`.
5. Merge only after required checks pass and review is complete. Delete merged branches when practical.

## GitHub repository settings

Branch protection cannot be enabled from this local checkout. In GitHub, open **Settings → Branches → Add branch protection rule** for `main` and configure:

- Require a pull request before merging; require one approval when collaborators are available.
- Require status checks to pass before merging; select **Backend CI / backend** and **Mobile CI / mobile**.
- Require branches to be up to date before merging.
- Block force pushes and branch deletion.
- Restrict direct pushes to `main` where the plan permits.

The workflows also run on pull requests targeting `main`. Do not enable production deployment from pull request events.
