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
        </div>
      </div>
    </nav>
  </div>
  <div class="hero-body">
    <div class="container has-text-centered">
      <h1 class="title">SignOff</h1>
      <div class="field">
        <p class="control">
          <input class="input" type="text" value="{{ .name }}" disabled>
	</p>
      </div>
      <div class="field">
        <p class="control">
          <button class="button is-danger" id="delete">Delete User</button>
        </p>
      </div>
      <div>Can you really delete your account? All your notes will also be deleted.</div>
    </div>
  </div>
</section>