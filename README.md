<br />
<div align="center">

<h1 align="center">migrate-tool</h1>

  <p align="center">
    CLI tool to manage SQL migrations
    <br />
    <a href="https://github.com/Rex-82/migrate-tool/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    ·
    <a href="https://github.com/Rex-82/migrate-tool/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>
</div>

<div align="center">
  <img src="https://github.com/Rex-82/migrate-tool/blob/main/assets/sample.gif"/>
</div>

<h2>About</h2>

This tool comes with sensible defaults for test projects, such as:
- host: `localhost`
- port: `3306`
- username: `root`
- migrations dir: `./db/migrations`

### Generated files
Migration files follow the format `YYYYMMDD_hhmmss_db_type.sql`
- `20240920_174815_db_schema.sql`
- `20240920_174815_db_data.sql`
- `20240920_174815_db_full.sql`

<h2>How to use</h2>

1. <a href="https://github.com/Rex-82/migrate-tool/releases">Download the latest release</a> that matches your architecture and OS
2. Extract the archive
3. _(suggested)_ Move the executable to one of your PATH dirs
4. Run the executable from within your project

## License

Distributed under the MIT License. See <a href="https://github.com/Rex-82/migrate-tool/blob/main/LICENSE">LICENSE</a> for more information.
