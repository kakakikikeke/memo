<section class="hero is-primary is-fullheight">
  <div class="hero-head">
    <nav class="navbar" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <a role="button" class="navbar-burger" aria-label="menu" aria-expanded="false" data-target="memo-navbar">
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </a>
      </div>
      <div id="memo-navbar" class="navbar-menu">
        <div class="navbar-start">
          <div class="navbar-item">
            <a class="button is-primary" href="/">
              <span>Home</span>
            </a>
          </div>
          <div class="navbar-item">
            <a class="button is-primary" href="/image">
              <span>Board</span>
            </a>
          </div>
          <div class="navbar-item">
            <a class="button is-primary" href="/file">
              <span>File</span>
            </a>
          </div>
          <div class="navbar-item">
            <div class="buttons">
              <a class="button is-primary is-inverted" href="https://github.com/kakakikikeke/memo" target="_blank">
                <span class="icon">
                  <i class="fa fa-github"></i>
                </span>
                <span>Show code</span>
              </a>
            </div>
          </div>
        </div>
        <div class="navbar-end">
          {{ if .isLogin }}
          <div class="navbar-item">
            <a class="button is-primary" id="logout" href="/logout">
              <span>Logout ({{ .name }})</span>
            </a>
          </div>
          <div class="navbar-item">
            <a class="button is-danger" id="signoff" href="/signoff">
              <span>Delete user ({{ .name }})</span>
            </a>
          </div>
          {{ else }}
          <div class="navbar-item">
            <a class="button is-primary" href="/login">
              <span>Login</span>
            </a>
          </div>
          <div class="navbar-item">
            <a class="button is-primary" href="/signup">
              <span>SignUp</span>
            </a>
          </div>
          {{ end }}
        </div>
      </div>
    </nav>
  </div>
  <div class="hero-body">
    <div class="container has-text-centered">
      <h1 class="title">file</h1>
      <h2 class="subtitle">
        Free files for Everyone.
      </h2>
      <div class="file has-name is-fullwidth">
        <label class="file-label">
          <input class="file-input" id="upload_file" type="file" name="file">
          <span class="file-cta">
            <span class="file-icon">
              <i class="fa fa-upload"></i>
            </span>
            <span class="file-label">
              Choose a file…
            </span>
          </span>
          <span class="file-name" id="file_name">
          </span>
        </label>
      </div>
      <div id="msg"></div>
      <div class="field mt-3">
        <p class="control">
          <button class="button is-info" id="upload" disabled>Upload</button>
        </p>
      </div>
      {{ range $file := .files }}
        {{ $content := get_content $file }}
        <div><a download="{{ get_file_name $file }}" {{ $content | attr }}>{{ get_file_name $file }}</a></div>
      {{ end }}
    </div>
  </div>
  <div class="hero-foot">
    <nav class="tabs is-centered">
      <div class="container">
        <a class="button is-danger" id="clear_file">Clear</a>
      </div>
    </nav>
  </div>
</section>