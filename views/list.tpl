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
      <h1 class="title">memo</h1>
      <h2 class="subtitle">
        Free note for Everyone.
      </h2>
      <div class="field">
        <p class="control">
          <label class="checkbox">
            <input type="checkbox" id="toggle">Textbox</input>
          </label>
        </p>
      </div>
      <div class="field">
        <p class="control">
          <input class="input" type="text" placeholder="something" id="value">
	</p>
        <p class="control">
          <textarea class="textarea" placeholder="some lines" id="values" style="display:none"></textarea>
        </p>
      </div>
      <div class="field">
        <p class="control">
          <button class="button is-info" id="submit" disabled>memo</button>
        </p>
      </div>
      {{ range $memo := .memos }}
        {{ $replaced_memo := rep $memo "\n" "<br>" }}
        <div>{{ str2html $replaced_memo }}</div>
      {{ end }}
    </div>
  </div>
  <div class="hero-foot">
    <nav class="tabs is-centered">
      <div class="container">
        <a class="button is-danger" id="clear">Clear</a>
      </div>
    </nav>
  </div>
</section>