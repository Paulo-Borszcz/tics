<div id="page"><link rel="stylesheet" type="text/css" href="/public/lib/prismjs.min.css?v=5988251963bdbddef4ab6a4dcf079a377ada5049" media="all"><script type="text/javascript" src="/public/lib/prismjs.min.js?v=5988251963bdbddef4ab6a4dcf079a377ada5049"></script><div class="documentation"><h1 id="glpi-rest-api%3A--documentation">GLPI REST API:  Documentation</h1>

<h2 id="summary">Summary</h2>

<ul>
<li><a href="#glossary">Glossary</a></li>
<li><a href="#important">Important</a></li>
<li><a href="#init-session">Init session</a></li>
<li><a href="#kill-session">Kill session</a></li>
<li><a href="#lost-password">Lost password</a></li>
<li><a href="#get-my-profiles">Get my profiles</a></li>
<li><a href="#get-active-profile">Get active profile</a></li>
<li><a href="#change-active-profile">Change active profile</a></li>
<li><a href="#get-my-entities">Get my entities</a></li>
<li><a href="#get-active-entities">Get active entities</a></li>
<li><a href="#change-active-entities">Change active entities</a></li>
<li><a href="#get-full-session">Get full session</a></li>
<li><a href="#get-glpi-config">Get GLPI config</a></li>
<li><a href="#get-an-item">Get an item</a></li>
<li><a href="#get-all-items">Get all items</a></li>
<li><a href="#get-sub-items">Get sub items</a></li>
<li><a href="#get-multiple-items">Get multiple items</a></li>
<li><a href="#list-searchoptions">List searchOptions</a></li>
<li><a href="#search-items">Search items</a></li>
<li><a href="#add-items">Add item(s)</a></li>
<li><a href="#update-items">Update item(s)</a></li>
<li><a href="#delete-items">Delete item(s)</a></li>
<li><a href="#get-available-massive-actions-for-an-itemtype">Get available massive actions for an itemtype</a></li>
<li><a href="#get-available-massive-actions-for-an-item">Get available massive actions for an item</a></li>
<li><a href="#get-massive-action-parameters">Get massive action parameters</a></li>
<li><a href="#apply-massive-action">Apply massive action</a></li>
<li><a href="#special-cases">Special cases</a></li>
<li><a href="#errors">Errors</a></li>
<li><a href="#servers-configuration">Servers configuration</a></li>
</ul>

<h2 id="glossary">Glossary</h2>

<dl>
<dt>Endpoint</dt>
<dd>Resource available though the API.
The endpoint is the URL where your API can be accessed by a client application</dd>

<dt>Method</dt>
<dd>HTTP verbs to indicate the desired action to be performed on the identified resource.
See: https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol#Request_methods</dd>

<dt>itemtype</dt>
<dd>A GLPI type, could be an asset, an ITIL or a configuration object, etc.
This type must be a class who inherits CommonDTBM GLPI class.
See <a href="https://forge.glpi-project.org/apidoc/class-CommonDBTM.html">List itemtypes</a>.</dd>

<dt>searchOption</dt>
<dd>A column identifier (integer) of an itemtype (ex: 1 -&gt; id, 2 -&gt; name, ...).
See <a href="#list-searchoptions">List searchOptions</a> endpoint.</dd>

<dt>JSON Payload</dt>
<dd>content of HTTP Request in JSON format (HTTP body)</dd>

<dt>Query string</dt>
<dd>URL parameters</dd>

<dt>User token</dt>
<dd>Used in login process instead of login/password couple.
It represent the user with a string.
You can find user token in the settings tabs of users.</dd>

<dt>Session token</dt>
<dd>A string describing a valid session in glpi.
Except initSession endpoint who provide this token, all others require this string to be used.</dd>

<dt>App(lication) token</dt>
<dd>An optional way to filter the access to the API.
On API call, it will try to find an API client matching your IP and the app token (if provided).
You can define an API client with an app token in general configuration for each of your external applications to identify them (each API client have its own history).</dd>
</dl>

<h2 id="important">Important</h2>

<ul>
<li><p>You should always provide a Content-Type header in your HTTP calls.
Currently, the API supports:</p>

<ul>
<li>application/json</li>
<li>multipart/form-data (for files upload, see <a href="#add-item-s">Add item(s)</a> endpoint.</li>
</ul></li>
<li><p>GET requests must have an empty body. You must pass all parameters in URL.
Failing to do so will trigger an HTTP 400 response.</p></li>
<li><p>By default, sessions used in this API are read-only.
Only some methods have write access to session:</p>

<ul>
<li><a href="#init-session">initSession</a></li>
<li><a href="#kill-session">killSession</a></li>
<li><a href="#change-active-entities">changeActiveEntities</a></li>
<li><a href="#change-active-profile">changeActiveProfile</a></li>
</ul>

<p>You could pass an additional parameter "session_write=true" to bypass this default.
This read-only mode allow to use this API with parallel calls.
In write mode, sessions are locked and your client must wait the end of a call before the next one can execute.</p></li>
<li><p>You can filter API access by enable the following parameters in GLPI General Configuration (API tab):</p>

<ul>
<li>IPv4 range</li>
<li>IPv6 address</li>
<li><em>App-Token</em> parameter: if not empty, you should pass this parameter in all of your API calls</li>
</ul></li>
<li><p>Session and App tokens could be provided in query string instead of header parameters.</p></li>
</ul>

<h2 id="init-session">Init session</h2>

<ul>
<li><strong>URL</strong>: apirest.php/initSession/</li>
<li><strong>Description</strong>: Request a session token to uses other API endpoints.</li>
<li><strong>Method</strong>: GET</li>
<li><p><strong>Parameters</strong>: (Headers)</p>

<ul>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
<li>a couple <em>login</em> &amp; <em>password</em>: 2 parameters to login with user authentication.
You should pass this 2 parameters in <a href="https://en.wikipedia.org/wiki/Basic_access_authentication">http basic auth</a>.
It consists in a Base64 string with login and password separated by ":"
A valid Authorization header is:

<ul>
<li>"Authorization: Basic base64({login}:{password})"</li>
</ul></li>
</ul>

<blockquote>
  <p><strong>OR</strong></p>
</blockquote>

<ul>
<li>an <em>user_token</em> defined in User Preference (See 'Remote access key')
You should pass this parameter in 'Authorization' HTTP header.
A valid Authorization header is:

<ul>
<li>"Authorization: user_token q56hqkniwot8wntb3z1qarka5atf365taaa2uyjrn"</li>
</ul></li>
</ul></li>
<li><strong>Parameters</strong>: (query string)

<ul>
<li><em>get_full_session</em> (default: false): Get the full session, useful if you want to login and access session data in one request.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) with the <em>session_token</em> string.</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
<li>401 (UNAUTHORIZED)</li>
</ul></li>
</ul>

<p>Examples usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Authorization: Basic Z2xwaTpnbHBp"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/initSession'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> <span class="token punctuation">{</span>
   <span class="token string">"session_token"</span><span class="token builtin class-name">:</span> <span class="token string">"83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span>
<span class="token punctuation">}</span>

$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Authorization: user_token q56hqkniwot8wntb3z1qarka5atf365taaa2uyjrn"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/initSession?get_full_session=true'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> <span class="token punctuation">{</span>
   <span class="token string">"session_token"</span><span class="token builtin class-name">:</span> <span class="token string">"83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span>,
   <span class="token string">"session"</span><span class="token builtin class-name">:</span> <span class="token punctuation">{</span>
      <span class="token string">'glpi_plugins'</span><span class="token builtin class-name">:</span> <span class="token punctuation">..</span>.,
      <span class="token string">'glpicookietest'</span><span class="token builtin class-name">:</span> <span class="token punctuation">..</span>.,
      <span class="token string">'glpicsrftokens'</span><span class="token builtin class-name">:</span> <span class="token punctuation">..</span>.,
      <span class="token punctuation">..</span>.
   <span class="token punctuation">}</span>
<span class="token punctuation">}</span>
</code></pre>

<h2 id="kill-session">Kill session</h2>

<ul>
<li><strong>URL</strong>: apirest.php/killSession/</li>
<li><strong>Description</strong>: Destroy a session identified by a session token.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK).</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/killSession'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
</code></pre>

<h2 id="lost-password">Lost password</h2>

<p>This endpoint allows to request password recovery and password reset. This endpoint works under the following
conditions:
* GLPI has notifications enabled
* the email address of the user belongs to a user account.</p>

<p>Reset password request:</p>

<ul>
<li><strong>URL</strong>: apirest.php/lostPassword/</li>
<li><strong>Description</strong>: Sends a notification to the user to reset his password</li>
<li><strong>Method</strong>: PUT or PATCH</li>
<li><strong>Parameters</strong>: (JSON Payload)

<ul>
<li><em>email</em>: email address of the user to recover. Mandatory.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK).</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
</ul></li>
</ul>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X PUT <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"email": "user@domain.com"}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/lostPassword'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
</code></pre>

<p>Password reset :</p>

<ul>
<li><strong>URL</strong>: apirest.php/lostPassword/</li>
<li><strong>Description</strong>: Sends a notification to the user to reset his password</li>
<li><strong>Method</strong>: PUT or PATCH</li>
<li><strong>Parameters</strong>: (JSON Payload)

<ul>
<li><em>email</em>: email address of the user to recover. Mandatory.</li>
<li><em>password_forget_token</em>: reset token</li>
<li><em>password</em>: the new password for the user</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK).</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
</ul></li>
</ul>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X PUT <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"email": "user@domain.com", \
     "password_forget_token": "b0a4cfe81448299ebed57442f4f21929c80ebee5" \
     "password": "NewPassword" \
    }'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/lostPassword'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
</code></pre>

<h2 id="get-my-profiles">Get my profiles</h2>

<ul>
<li><strong>URL</strong>: <a href="getMyProfiles/?debug">apirest.php/getMyProfiles/</a></li>
<li><strong>Description</strong>: Return all the profiles associated to logged user.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) with an array of all profiles.</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getMyProfiles'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> <span class="token punctuation">{</span>
   <span class="token string">'myprofiles'</span><span class="token builtin class-name">:</span> <span class="token punctuation">[</span>
      <span class="token punctuation">{</span>
         <span class="token string">'id'</span><span class="token builtin class-name">:</span> <span class="token number">1</span>
         <span class="token string">'name'</span><span class="token builtin class-name">:</span> <span class="token string">"Super-admin"</span>,
         <span class="token string">'entities'</span><span class="token builtin class-name">:</span> <span class="token punctuation">[</span>
            <span class="token punctuation">..</span>.
         <span class="token punctuation">]</span>,
         <span class="token punctuation">..</span>.
      <span class="token punctuation">}</span>,
      <span class="token punctuation">..</span><span class="token punctuation">..</span>
   <span class="token punctuation">]</span>
</code></pre>

<h2 id="get-active-profile">Get active profile</h2>

<ul>
<li><strong>URL</strong>: <a href="getActiveProfile/?debug">apirest.php/getActiveProfile/</a></li>
<li><strong>Description</strong>: Return the current active profile.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) with an array representing current profile.</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getActiveProfile'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> <span class="token punctuation">{</span>
      <span class="token string">'name'</span><span class="token builtin class-name">:</span> <span class="token string">"Super-admin"</span>,
      <span class="token string">'entities'</span><span class="token builtin class-name">:</span> <span class="token punctuation">[</span>
         <span class="token punctuation">..</span>.
      <span class="token punctuation">]</span>
   <span class="token punctuation">}</span>
</code></pre>

<h2 id="change-active-profile">Change active profile</h2>

<ul>
<li><strong>URL</strong>: <a href="changeActiveProfile/?profiles_id=4&amp;debug">apirest.php/changeActiveProfile/</a></li>
<li><strong>Description</strong>: Change active profile to the profiles_id one. See <a href="#get-my-profiles">getMyProfiles</a> endpoint for possible profiles.</li>
<li><strong>Method</strong>: POST</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (JSON Payload)

<ul>
<li><em>profiles_id</em>: (default 'all') ID of the new active profile. Mandatory.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK).</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
<li>404 (Not found) with a message indicating an error ig the profile does not exists or usable.</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X POST <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"profiles_id": 4}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/changeActiveProfile'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
</code></pre>

<h2 id="get-my-entities">Get my entities</h2>

<ul>
<li><strong>URL</strong>: <a href="getMyEntities/?debug">apirest.php/getMyEntities/</a></li>
<li><strong>Description</strong>: Return all the possible entities of the current logged user (and for current active profile).</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (query string)

<ul>
<li><em>is_recursive</em> (default: false): Also display sub entities of the active entity. Optionnal</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) with an array of all entities (with id and name).</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getMyEntities'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> <span class="token punctuation">{</span>
   <span class="token string">'myentities'</span><span class="token builtin class-name">:</span> <span class="token punctuation">[</span>
     <span class="token punctuation">{</span>
       <span class="token string">'id'</span><span class="token builtin class-name">:</span>   <span class="token number">71</span>
       <span class="token string">'name'</span><span class="token builtin class-name">:</span> <span class="token string">"my_entity"</span>
     <span class="token punctuation">}</span>,
   <span class="token punctuation">..</span><span class="token punctuation">..</span>
   <span class="token punctuation">]</span>
  <span class="token punctuation">}</span>
</code></pre>

<h2 id="get-active-entities">Get active entities</h2>

<ul>
<li><strong>URL</strong>: <a href="getActiveEntities/?debug">apirest.php/getActiveEntities/</a></li>
<li><strong>Description</strong>: Return active entities of current logged user.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) with an array with 3 keys:</li>
<li><em>active_entity</em>: current set entity.</li>
<li><em>active_entity_recursive</em>: boolean, if we see sons of this entity.</li>
<li><em>active_entities</em>: array all active entities (active_entity and its sons).</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getActiveEntities'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> <span class="token punctuation">{</span>
   <span class="token string">'active_entity'</span><span class="token builtin class-name">:</span> <span class="token punctuation">{</span>
      <span class="token string">'id'</span><span class="token builtin class-name">:</span> <span class="token number">1</span>,
      <span class="token string">'active_entity_recursive'</span><span class="token builtin class-name">:</span> true,
      <span class="token string">'active_entities'</span><span class="token builtin class-name">:</span> <span class="token punctuation">[</span>
        <span class="token punctuation">{</span><span class="token string">"id"</span>:1<span class="token punctuation">}</span>,
        <span class="token punctuation">{</span><span class="token string">"id"</span>:71<span class="token punctuation">}</span>,<span class="token punctuation">..</span>.
      <span class="token punctuation">]</span>
   <span class="token punctuation">}</span>
<span class="token punctuation">}</span>
</code></pre>

<h2 id="change-active-entities">Change active entities</h2>

<ul>
<li><strong>URL</strong>: <a href="changeActiveEntities/?entities_id=1&amp;is_recursive=0&amp;debug">apirest.php/changeActiveEntities/</a></li>
<li><strong>Description</strong>: Change active entity to the entities_id one. See <a href="#get-my-entities">getMyEntities</a> endpoint for possible entities.</li>
<li><strong>Method</strong>: POST</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (JSON Payload)

<ul>
<li><em>entities_id</em>: (default 'all') ID of the new active entity ("all" =&gt; load all possible entities). Optional.</li>
<li><em>is_recursive</em>: (default false) Also display sub entities of the active entity.  Optional.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK).</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X POST <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"entities_id": 1, "is_recursive": true}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/changeActiveEntities'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
</code></pre>

<h2 id="get-full-session">Get full session</h2>

<ul>
<li><strong>URL</strong>: <a href="getFullSession/?debug">apirest.php/getFullSession/</a></li>
<li><strong>Description</strong>: Return the current php $_SESSION.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) with an array representing the php session.</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getFullSession'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> <span class="token punctuation">{</span>
      <span class="token string">'session'</span><span class="token builtin class-name">:</span> <span class="token punctuation">{</span>
         <span class="token string">'glpi_plugins'</span><span class="token builtin class-name">:</span> <span class="token punctuation">..</span>.,
         <span class="token string">'glpicookietest'</span><span class="token builtin class-name">:</span> <span class="token punctuation">..</span>.,
         <span class="token string">'glpicsrftokens'</span><span class="token builtin class-name">:</span> <span class="token punctuation">..</span>.,
         <span class="token punctuation">..</span>.
      <span class="token punctuation">}</span>
   <span class="token punctuation">}</span>
</code></pre>

<h2 id="get-glpi-config">Get GLPI config</h2>

<ul>
<li><strong>URL</strong>: <a href="getGlpiConfig/?debug">apirest.php/getGlpiConfig/</a></li>
<li><strong>Description</strong>: Return the current $CFG_GLPI.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) with an array representing the php global variable $CFG_GLPI.</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getGlpiConfig'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> <span class="token punctuation">{</span>
      <span class="token string">'cfg_glpi'</span><span class="token builtin class-name">:</span> <span class="token punctuation">{</span>
         <span class="token string">'languages'</span><span class="token builtin class-name">:</span> <span class="token punctuation">..</span>.,
         <span class="token string">'glpitables'</span><span class="token builtin class-name">:</span> <span class="token punctuation">..</span>.,
         <span class="token string">'unicity_types'</span>:<span class="token punctuation">..</span>.,
         <span class="token punctuation">..</span>.
      <span class="token punctuation">}</span>
   <span class="token punctuation">}</span>
</code></pre>

<h2 id="get-an-item">Get an item</h2>

<ul>
<li><strong>URL</strong>: <a href="User/2?debug">apirest.php/:itemtype/:id</a></li>
<li><strong>Description</strong>: Return the instance fields of itemtype identified by id.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (query string)

<ul>
<li><em>id</em>: unique identifier of the itemtype. Mandatory.</li>
<li><em>expand_dropdowns</em> (default: false): show dropdown name instead of id. Optional.</li>
<li><em>get_hateoas</em> (default: true): Show relations of the item in a links attribute. Optional.</li>
<li><em>get_sha1</em> (default: false): Get a sha1 signature instead of the full answer. Optional.</li>
<li><em>with_devices</em>: Only for [Computer, NetworkEquipment, Peripheral, Phone, Printer], retrieve the associated components. Optional.</li>
<li><em>with_disks</em>: Only for Computer, retrieve the associated file-systems. Optional.</li>
<li><em>with_softwares</em>: Only for Computer, retrieve the associated software's installations. Optional.</li>
<li><em>with_connections</em>: Only for Computer, retrieve the associated direct connections (like peripherals and printers) .Optional.</li>
<li><em>with_networkports</em>: Retrieve all network connections and advanced information. Optionnal.</li>
<li><em>with_infocoms</em>: Retrieve financial and administrative information. Optional.</li>
<li><em>with_contracts</em>: Retrieve associated contracts. Optional.</li>
<li><em>with_documents</em>: Retrieve associated external documents. Optional.</li>
<li><em>with_tickets</em>: Retrieve associated ITIL tickets. Optional.</li>
<li><em>with_problems</em>: Retrieve associated ITIL problems. Optional.</li>
<li><em>with_changes</em>: Retrieve associated ITIL changes. Optional.</li>
<li><em>with_notes</em>: Retrieve Notes. Optional.</li>
<li><em>with_logs</em>: Retrieve historical. Optional.</li>
<li><em>add_keys_names</em>: Retrieve friendly names. Array containing fkey(s) and/or "id". Optional.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) with item data (Last-Modified header should contain the date of last modification of the item).</li>
<li>401 (UNAUTHORIZED).</li>
<li>404 (NOT FOUND).</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Computer/71?expand_dropdowns=true'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> <span class="token punctuation">{</span>
    <span class="token string">"id"</span><span class="token builtin class-name">:</span> <span class="token number">71</span>,
    <span class="token string">"entities_id"</span><span class="token builtin class-name">:</span> <span class="token string">"Root Entity"</span>,
    <span class="token string">"name"</span><span class="token builtin class-name">:</span> <span class="token string">"adelaunay-ThinkPad-Edge-E320"</span>,
    <span class="token string">"serial"</span><span class="token builtin class-name">:</span> <span class="token string">"12345"</span>,
    <span class="token string">"otherserial"</span><span class="token builtin class-name">:</span> <span class="token string">"test2"</span>,
    <span class="token string">"contact"</span><span class="token builtin class-name">:</span> <span class="token string">"adelaunay"</span>,
    <span class="token string">"contact_num"</span><span class="token builtin class-name">:</span> null,
    <span class="token string">"users_id_tech"</span><span class="token builtin class-name">:</span> <span class="token string">" "</span>,
    <span class="token string">"groups_id_tech"</span><span class="token builtin class-name">:</span> <span class="token string">" "</span>,
    <span class="token string">"comment"</span><span class="token builtin class-name">:</span> <span class="token string">"test222222qsdqsd"</span>,
    <span class="token string">"date_mod"</span><span class="token builtin class-name">:</span> <span class="token string">"2015-09-25 09:33:41"</span>,
    <span class="token string">"autoupdatesystems_id"</span><span class="token builtin class-name">:</span> <span class="token string">" "</span>,
    <span class="token string">"locations_id"</span><span class="token builtin class-name">:</span> <span class="token string">"00:0e:08:3b:7d:04"</span>,
    <span class="token string">"networks_id"</span><span class="token builtin class-name">:</span> <span class="token string">" "</span>,
    <span class="token string">"computermodels_id"</span><span class="token builtin class-name">:</span> <span class="token string">"1298A8G"</span>,
    <span class="token string">"computertypes_id"</span><span class="token builtin class-name">:</span> <span class="token string">"Notebook"</span>,
    <span class="token string">"is_template"</span><span class="token builtin class-name">:</span> <span class="token number">0</span>,
    <span class="token string">"template_name"</span><span class="token builtin class-name">:</span> null,
    <span class="token string">"manufacturers_id"</span><span class="token builtin class-name">:</span> <span class="token string">"LENOVO"</span>,
    <span class="token string">"is_deleted"</span><span class="token builtin class-name">:</span> <span class="token number">0</span>,
    <span class="token string">"is_dynamic"</span><span class="token builtin class-name">:</span> <span class="token number">1</span>,
    <span class="token string">"users_id"</span><span class="token builtin class-name">:</span> <span class="token string">"adelaunay"</span>,
    <span class="token string">"groups_id"</span><span class="token builtin class-name">:</span> <span class="token string">" "</span>,
    <span class="token string">"states_id"</span><span class="token builtin class-name">:</span> <span class="token string">"Production"</span>,
    <span class="token string">"ticket_tco"</span><span class="token builtin class-name">:</span> <span class="token string">"0.0000"</span>,
    <span class="token string">"uuid"</span><span class="token builtin class-name">:</span> <span class="token string">""</span>,
    <span class="token string">"date_creation"</span><span class="token builtin class-name">:</span> null,
    <span class="token string">"links"</span><span class="token builtin class-name">:</span> <span class="token punctuation">[</span><span class="token punctuation">{</span>
       <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"Entity"</span>,
       <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/Entity/0"</span>
    <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
       <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"Location"</span>,
       <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/Location/3"</span>
    <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
       <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"Domain"</span>,
       <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/Domain/18"</span>
    <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
       <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"ComputerModel"</span>,
       <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/ComputerModel/11"</span>
    <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
       <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"ComputerType"</span>,
       <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/ComputerType/3"</span>
    <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
       <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"Manufacturer"</span>,
       <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/Manufacturer/260"</span>
    <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
       <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"User"</span>,
       <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/User/27"</span>
    <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
       <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"State"</span>,
       <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/State/1"</span>
    <span class="token punctuation">}</span><span class="token punctuation">]</span>
<span class="token punctuation">}</span>
</code></pre>

<p>Note: To download a document see <a href="#download-a-document-file">Download a document file</a>.</p>

<h2 id="get-all-items">Get all items</h2>

<ul>
<li><strong>URL</strong>: <a href="Computer/?debug">apirest.php/:itemtype/</a></li>
<li><strong>Description</strong>: Return a collection of rows of the itemtype.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (query string)

<ul>
<li><em>expand_dropdowns</em> (default: false): show dropdown name instead of id. Optional.</li>
<li><em>get_hateoas</em> (default: true): Show relation of item in a links attribute. Optional.</li>
<li><em>only_id</em> (default: false): keep only id keys in returned data. Optional.</li>
<li><em>range</em> (default: 0-49):  a string with a couple of number for start and end of pagination separated by a '-'. Ex: 150-199. Optional.</li>
<li><em>sort</em> (default 1): name of the field to sort by. Optional.</li>
<li><em>order</em> (default ASC): ASC - Ascending sort / DESC Descending sort. Optional.</li>
<li><em>searchText</em> (default NULL): array of filters to pass on the query (with key = field and value the text to search)</li>
<li><em>is_deleted</em> (default: false): Return deleted element. Optional.</li>
<li><em>add_keys_names</em>: Retrieve friendly names. Array containing fkey(s) and/or "id". Optional.</li>
<li><em>with_networkports</em>: Retrieve all network connections and advanced information. Optionnal.</li>
</ul></li>
<li><p><strong>Returns</strong>:</p>

<ul>
<li>200 (OK) with items data.</li>
<li>206 (PARTIAL CONTENT) with items data defined by range.</li>
<li>401 (UNAUTHORIZED).</li>
</ul>

<p>and theses headers:</p>

<ul>
<li><em>Content-Range</em> offset – limit / count</li>
<li><em>Accept-Range</em> itemtype max</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Computer/?expand_dropdowns=true'</span>

<span class="token operator">&lt;</span> <span class="token number">206</span> OK
<span class="token operator">&lt;</span> Content-Range: <span class="token number">0</span>-49/200
<span class="token operator">&lt;</span> Accept-Range: <span class="token number">990</span>
<span class="token operator">&lt;</span> <span class="token punctuation">[</span>
   <span class="token punctuation">{</span>
      <span class="token string">"id"</span><span class="token builtin class-name">:</span> <span class="token number">34</span>,
      <span class="token string">"entities_id"</span><span class="token builtin class-name">:</span> <span class="token string">"Root Entity"</span>,
      <span class="token string">"name"</span><span class="token builtin class-name">:</span> <span class="token string">"glpi"</span>,
      <span class="token string">"serial"</span><span class="token builtin class-name">:</span> <span class="token string">"VMware-42 01 f4 65 27 59 a9 fb-11 bc cd b8 64 68 1f 4b"</span>,
      <span class="token string">"otherserial"</span><span class="token builtin class-name">:</span> null,
      <span class="token string">"contact"</span><span class="token builtin class-name">:</span> <span class="token string">"teclib"</span>,
      <span class="token string">"contact_num"</span><span class="token builtin class-name">:</span> null,
      <span class="token string">"users_id_tech"</span><span class="token builtin class-name">:</span> <span class="token string">"&amp;nbsp;"</span>,
      <span class="token string">"groups_id_tech"</span><span class="token builtin class-name">:</span> <span class="token string">"&amp;nbsp;"</span>,
      <span class="token string">"comment"</span><span class="token builtin class-name">:</span> <span class="token string">"x86_64/00-09-15 08:03:28"</span>,
      <span class="token string">"date_mod"</span><span class="token builtin class-name">:</span> <span class="token string">"2011-12-16 17:52:55"</span>,
      <span class="token string">"autoupdatesystems_id"</span><span class="token builtin class-name">:</span> <span class="token string">"FusionInventory"</span>,
      <span class="token string">"locations_id"</span><span class="token builtin class-name">:</span> <span class="token string">"&amp;nbsp;"</span>,
      <span class="token string">"networks_id"</span><span class="token builtin class-name">:</span> <span class="token string">"&amp;nbsp;"</span>,
      <span class="token string">"computermodels_id"</span><span class="token builtin class-name">:</span> <span class="token string">"VMware Virtual Platform"</span>,
      <span class="token string">"computertypes_id"</span><span class="token builtin class-name">:</span> <span class="token string">"Other"</span>,
      <span class="token string">"is_template"</span><span class="token builtin class-name">:</span> <span class="token number">0</span>,
      <span class="token string">"template_name"</span><span class="token builtin class-name">:</span> null,
      <span class="token string">"manufacturers_id"</span><span class="token builtin class-name">:</span> <span class="token string">"VMware, Inc."</span>,
      <span class="token string">"is_deleted"</span><span class="token builtin class-name">:</span> <span class="token number">0</span>,
      <span class="token string">"is_dynamic"</span><span class="token builtin class-name">:</span> <span class="token number">1</span>,
      <span class="token string">"users_id"</span><span class="token builtin class-name">:</span> <span class="token string">"&amp;nbsp;"</span>,
      <span class="token string">"groups_id"</span><span class="token builtin class-name">:</span> <span class="token string">"&amp;nbsp;"</span>,
      <span class="token string">"states_id"</span><span class="token builtin class-name">:</span> <span class="token string">"Production"</span>,
      <span class="token string">"ticket_tco"</span><span class="token builtin class-name">:</span> <span class="token string">"0.0000"</span>,
      <span class="token string">"uuid"</span><span class="token builtin class-name">:</span> <span class="token string">"4201F465-2759-A9FB-11BC-CDB864681F4B"</span>,
      <span class="token string">"links"</span><span class="token builtin class-name">:</span> <span class="token punctuation">[</span><span class="token punctuation">{</span>
         <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"Entity"</span>,
         <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/Entity/0"</span>
      <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
         <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"AutoUpdateSystem"</span>,
         <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/AutoUpdateSystem/1"</span>
      <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
         <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"Domain"</span>,
         <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/Domain/12"</span>
      <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
         <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"ComputerModel"</span>,
         <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/ComputerModel/1"</span>
      <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
         <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"ComputerType"</span>,
         <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/ComputerType/2"</span>
      <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
         <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"Manufacturer"</span>,
         <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/Manufacturer/1"</span>
      <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
         <span class="token string">"rel"</span><span class="token builtin class-name">:</span> <span class="token string">"State"</span>,
         <span class="token string">"href"</span><span class="token builtin class-name">:</span> <span class="token string">"http://path/to/glpi/api/State/1"</span>
      <span class="token punctuation">}</span><span class="token punctuation">]</span>
   <span class="token punctuation">}</span>,
   <span class="token punctuation">{</span>
      <span class="token string">"id"</span><span class="token builtin class-name">:</span> <span class="token number">35</span>,
      <span class="token string">"entities_id"</span><span class="token builtin class-name">:</span> <span class="token string">"Root Entity"</span>,
      <span class="token string">"name"</span><span class="token builtin class-name">:</span> <span class="token string">"mavm1"</span>,
      <span class="token string">"serial"</span><span class="token builtin class-name">:</span> <span class="token string">"VMware-42 20 d3 04 ac 49 ed c8-ea 15 50 49 e1 40 0f 6c"</span>,
      <span class="token string">"otherserial"</span><span class="token builtin class-name">:</span> null,
      <span class="token string">"contact"</span><span class="token builtin class-name">:</span> <span class="token string">"teclib"</span>,
      <span class="token string">"contact_num"</span><span class="token builtin class-name">:</span> null,
      <span class="token string">"users_id_tech"</span><span class="token builtin class-name">:</span> <span class="token string">"&amp;nbsp;"</span>,
      <span class="token string">"groups_id_tech"</span><span class="token builtin class-name">:</span> <span class="token string">"&amp;nbsp;"</span>,
      <span class="token string">"comment"</span><span class="token builtin class-name">:</span> <span class="token string">"x86_64/01-01-04 19:50:40"</span>,
      <span class="token string">"date_mod"</span><span class="token builtin class-name">:</span> <span class="token string">"2012-05-24 06:43:35"</span>,
      <span class="token punctuation">..</span>.
   <span class="token punctuation">}</span>
<span class="token punctuation">]</span>
</code></pre>

<h2 id="get-sub-items">Get sub items</h2>

<ul>
<li><strong>URL</strong>: <a href="User/2/Log?debug">apirest.php/:itemtype/:id/:sub_itemtype</a></li>
<li><strong>Description</strong>: Return a collection of rows of the sub_itemtype for the identified item.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (query string)

<ul>
<li>id: unique identifier of the parent itemtype. Mandatory.</li>
<li><em>expand_dropdowns</em> (default: false): show dropdown name instead of id. Optional.</li>
<li><em>get_hateoas</em> (default: true): Show item's relations in a links attribute. Optional.</li>
<li><em>only_id</em> (default: false): keep only id keys in returned data. Optional.</li>
<li><em>range</em> (default: 0-49): a string with a couple of number for start and end of pagination separated by a '-' char. Ex: 150-199. Optional.</li>
<li><em>sort</em> (default 1): id of the "searchoption" to sort by. Optional.</li>
<li><em>order</em> (default ASC): ASC - Ascending sort / DESC Descending sort. Optional.</li>
<li><em>add_keys_names</em>: Retrieve friendly names. Array containing fkey(s) and/or "id". Optional.</li>
</ul></li>
<li><p><strong>Returns</strong>:</p>

<ul>
<li>200 (OK) with the items data.</li>
<li>401 (UNAUTHORIZED).</li>
</ul>

<p>and theses headers:</p>

<ul>
<li><em>Content-Range</em> offset – limit / count</li>
<li><em>Accept-Range</em> itemtype max</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/User/2/Log'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> Content-Range: <span class="token number">0</span>-49/200
<span class="token operator">&lt;</span> Accept-Range: <span class="token number">990</span>
<span class="token operator">&lt;</span> <span class="token punctuation">[</span>
   <span class="token punctuation">{</span>
      <span class="token string">"id"</span><span class="token builtin class-name">:</span> <span class="token number">22117</span>,
      <span class="token string">"itemtype"</span><span class="token builtin class-name">:</span> <span class="token string">"User"</span>,
      <span class="token string">"items_id"</span><span class="token builtin class-name">:</span> <span class="token number">2</span>,
      <span class="token string">"itemtype_link"</span><span class="token builtin class-name">:</span> <span class="token string">"Profile"</span>,
      <span class="token string">"linked_action"</span><span class="token builtin class-name">:</span> <span class="token number">17</span>,
      <span class="token string">"user_name"</span><span class="token builtin class-name">:</span> <span class="token string">"glpi (27)"</span>,
      <span class="token string">"date_mod"</span><span class="token builtin class-name">:</span> <span class="token string">"2015-10-13 10:00:59"</span>,
      <span class="token string">"id_search_option"</span><span class="token builtin class-name">:</span> <span class="token number">0</span>,
      <span class="token string">"old_value"</span><span class="token builtin class-name">:</span> <span class="token string">""</span>,
      <span class="token string">"new_value"</span><span class="token builtin class-name">:</span> <span class="token string">"super-admin (4)"</span>
   <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
      <span class="token string">"id"</span><span class="token builtin class-name">:</span> <span class="token number">22118</span>,
      <span class="token string">"itemtype"</span><span class="token builtin class-name">:</span> <span class="token string">"User"</span>,
      <span class="token string">"items_id"</span><span class="token builtin class-name">:</span> <span class="token number">2</span>,
      <span class="token string">"itemtype_link"</span><span class="token builtin class-name">:</span> <span class="token string">""</span>,
      <span class="token string">"linked_action"</span><span class="token builtin class-name">:</span> <span class="token number">0</span>,
      <span class="token string">"user_name"</span><span class="token builtin class-name">:</span> <span class="token string">"glpi (2)"</span>,
      <span class="token string">"date_mod"</span><span class="token builtin class-name">:</span> <span class="token string">"2015-10-13 10:01:22"</span>,
      <span class="token string">"id_search_option"</span><span class="token builtin class-name">:</span> <span class="token number">80</span>,
      <span class="token string">"old_value"</span><span class="token builtin class-name">:</span> <span class="token string">"Root entity (0)"</span>,
      <span class="token string">"new_value"</span><span class="token builtin class-name">:</span> <span class="token string">"Root entity &gt; my entity (1)"</span>
   <span class="token punctuation">}</span>, <span class="token punctuation">{</span>
      <span class="token punctuation">..</span>.
   <span class="token punctuation">}</span>
<span class="token punctuation">]</span>
</code></pre>

<h2 id="get-multiple-items">Get multiple items</h2>

<ul>
<li><strong>URL</strong>: apirest.php/getMultipleItems</li>
<li><strong>Description</strong>: Virtually call <a href="#get-an-item">Get an item</a> for each line in input. So, you can have a ticket, a user in the same query.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (query string)

<ul>
<li><em>items</em>: items to retrieve. Mandatory.
      Each line of this array should contains two keys:
      * itemtype
      * items_id</li>
<li><em>expand_dropdowns</em> (default: false): show dropdown name instead of id. Optional.</li>
<li><em>get_hateoas</em> (default: true): Show relations of the item in a links attribute. Optional.</li>
<li><em>get_sha1</em> (default: false): Get a sha1 signature instead of the full answer. Optional.</li>
<li><em>with_devices</em>: Only for [Computer, NetworkEquipment, Peripheral, Phone, Printer], retrieve the associated components. Optional.</li>
<li><em>with_disks</em>: Only for Computer, retrieve the associated file-systems. Optional.</li>
<li><em>with_softwares</em>: Only for Computer, retrieve the associated software's installations. Optional.</li>
<li><em>with_connections</em>: Only for Computer, retrieve the associated direct connections (like peripherals and printers) .Optional.</li>
<li><em>with_networkports</em>: Retrieve all network connections and advanced information. Optionnal.</li>
<li><em>with_infocoms</em>: Retrieve financial and administrative information. Optional.</li>
<li><em>with_contracts</em>: Retrieve associated contracts. Optional.</li>
<li><em>with_documents</em>: Retrieve associated external documents. Optional.</li>
<li><em>with_tickets</em>: Retrieve associated ITIL tickets. Optional.</li>
<li><em>with_problems</em>: Retrieve associated ITIL problems. Optional.</li>
<li><em>with_changes</em>: Retrieve associated ITIL changes. Optional.</li>
<li><em>with_notes</em>: Retrieve Notes. Optional.</li>
<li><em>with_logs</em>: Retrieve historical. Optional.</li>
<li><em>add_keys_names</em>: Retrieve friendly names. Array containing fkey(s) and/or "id". Optional.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) with item data (Last-Modified header should contain the date of last modification of the item).</li>
<li>401 (UNAUTHORIZED).</li>
<li>404 (NOT FOUND).</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"items": [{"itemtype": "User", "items_id": 2}, {"itemtype": "Entity", "items_id": 0}]}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getMultipleItems?items\[0\]\[itemtype\]\=User&amp;items\[0\]\[items_id\]\=2&amp;items\[1\]\[itemtype\]\=Entity&amp;items\[1\]\[items_id\]\=0'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> Content-Range: <span class="token number">0</span>-49/200
<span class="token operator">&lt;</span> Accept-Range: <span class="token number">990</span>
<span class="token operator">&lt;</span> <span class="token punctuation">[</span><span class="token punctuation">{</span>
   <span class="token string">"id"</span><span class="token builtin class-name">:</span> <span class="token number">2</span>,
   <span class="token string">"name"</span><span class="token builtin class-name">:</span> <span class="token string">"glpi"</span>,
   <span class="token punctuation">..</span>.
<span class="token punctuation">}</span>, <span class="token punctuation">{</span>
   <span class="token string">"id"</span><span class="token builtin class-name">:</span> <span class="token number">0</span>,
   <span class="token string">"name"</span><span class="token builtin class-name">:</span> <span class="token string">"Root Entity"</span>,
   <span class="token punctuation">..</span>.
<span class="token punctuation">}</span><span class="token punctuation">]</span>
</code></pre>

<h2 id="list-searchoptions">List searchOptions</h2>

<ul>
<li><strong>URL</strong>: <a href="listSearchOptions/Computer?debug">apirest.php/listSearchOptions/:itemtype</a></li>
<li><strong>Description</strong>: List the searchoptions of provided itemtype. To use with <a href="#search_items">Search items</a>.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (query string)

<ul>
<li><em>raw</em>: return searchoption uncleaned (as provided by core)</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) with all searchoptions of specified itemtype (format: searchoption_id: {option_content}).</li>
<li>401 (UNAUTHORIZED).</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/listSearchOptions/Computer'</span>
<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> <span class="token punctuation">{</span>
    <span class="token string">"common"</span><span class="token builtin class-name">:</span> <span class="token string">"Characteristics"</span>,

    <span class="token number">1</span>: <span class="token punctuation">{</span>
      <span class="token string">'name'</span><span class="token builtin class-name">:</span> <span class="token string">'Name'</span>
      <span class="token string">'table'</span><span class="token builtin class-name">:</span> <span class="token string">'glpi_computers'</span>
      <span class="token string">'field'</span><span class="token builtin class-name">:</span> <span class="token string">'name'</span>
      <span class="token string">'linkfield'</span><span class="token builtin class-name">:</span> <span class="token string">'name'</span>
      <span class="token string">'datatype'</span><span class="token builtin class-name">:</span> <span class="token string">'itemlink'</span>
      <span class="token string">'uid'</span><span class="token builtin class-name">:</span> <span class="token string">'Computer.name'</span>
   <span class="token punctuation">}</span>,
   <span class="token number">2</span>: <span class="token punctuation">{</span>
      <span class="token string">'name'</span><span class="token builtin class-name">:</span> <span class="token string">'ID'</span>
      <span class="token string">'table'</span><span class="token builtin class-name">:</span> <span class="token string">'glpi_computers'</span>
      <span class="token string">'field'</span><span class="token builtin class-name">:</span> <span class="token string">'id'</span>
      <span class="token string">'linkfield'</span><span class="token builtin class-name">:</span> <span class="token string">'id'</span>
      <span class="token string">'datatype'</span><span class="token builtin class-name">:</span> <span class="token string">'number'</span>
      <span class="token string">'uid'</span><span class="token builtin class-name">:</span> <span class="token string">'Computer.id'</span>
   <span class="token punctuation">}</span>,
   <span class="token number">3</span>: <span class="token punctuation">{</span>
      <span class="token string">'name'</span><span class="token builtin class-name">:</span> <span class="token string">'Location'</span>
      <span class="token string">'table'</span><span class="token builtin class-name">:</span> <span class="token string">'glpi_locations'</span>
      <span class="token string">'field'</span><span class="token builtin class-name">:</span> <span class="token string">'completename'</span>
      <span class="token string">'linkfield'</span><span class="token builtin class-name">:</span> <span class="token string">'locations_id'</span>
      <span class="token string">'datatype'</span><span class="token builtin class-name">:</span> <span class="token string">'dropdown'</span>
      <span class="token string">'uid'</span><span class="token builtin class-name">:</span> <span class="token string">'Computer.Location.completename'</span>
   <span class="token punctuation">}</span>,
   <span class="token punctuation">..</span>.
<span class="token punctuation">}</span>
</code></pre>

<h2 id="search-items">Search items</h2>

<ul>
<li><strong>URL</strong>: <a href="search/Computer/?debug">apirest.php/search/:itemtype/</a></li>
<li><strong>Description</strong>: Expose the GLPI searchEngine and combine criteria to retrieve a list of elements of specified itemtype.
&gt; Note: you can use 'AllAssets' itemtype to retrieve a combination of all asset's types.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><p><strong>Parameters</strong>: (query string)</p>

<ul>
<li><p><em>criteria</em>: array of criterion objects to filter search. Optional.
You can optionally precise <code>meta=true</code> to pass a searchoption of another itemtype (meta-criteria).
Each criterion object must provide at least:</p>

<ul>
<li><em>link</em>: (optional for 1st element) logical operator in [AND, OR, AND NOT, AND NOT].</li>
</ul>

<p>And you can pass a direct searchoption usage :</p>

<ul>
<li><em>field</em>: id of the searchoption.</li>
<li><em>meta</em>: boolean, is this criterion a meta one ?</li>
<li><em>itemtype</em>: for meta=true criterion, precise the itemtype to use.</li>
<li><em>searchtype</em>: type of search in [contains¹, equals², notequals², lessthan, morethan, under, notunder].</li>
<li><em>value</em>: the value to search.</li>
</ul>

<p>Or a list of sub-nodes with the key:</p>

<ul>
<li><em>criteria</em>: nested criteria inside this criteria.</li>
</ul>

<p>Ex:</p>

<pre class="language-json" tabindex="0"><code class="language-json">...
<span class="token property">"criteria"</span><span class="token operator">:</span>
 <span class="token punctuation">[</span>
    <span class="token punctuation">{</span>
       <span class="token property">"field"</span><span class="token operator">:</span>      <span class="token number">1</span><span class="token punctuation">,</span>
       <span class="token property">"searchtype"</span><span class="token operator">:</span> 'contains'<span class="token punctuation">,</span>
       <span class="token property">"value"</span><span class="token operator">:</span>      ''
    <span class="token punctuation">}</span><span class="token punctuation">,</span> <span class="token punctuation">{</span>
       <span class="token property">"link"</span><span class="token operator">:</span>       'AND'<span class="token punctuation">,</span>
       <span class="token property">"field"</span><span class="token operator">:</span>      <span class="token number">31</span><span class="token punctuation">,</span>
       <span class="token property">"searchtype"</span><span class="token operator">:</span> 'equals'<span class="token punctuation">,</span>
       <span class="token property">"value"</span><span class="token operator">:</span>      <span class="token number">1</span>
    <span class="token punctuation">}</span><span class="token punctuation">,</span> <span class="token punctuation">{</span>
       <span class="token property">"link"</span><span class="token operator">:</span>       'AND'<span class="token punctuation">,</span>
       <span class="token property">"meta"</span><span class="token operator">:</span>       <span class="token boolean">true</span><span class="token punctuation">,</span>
       <span class="token property">"itemtype"</span><span class="token operator">:</span>   'User'<span class="token punctuation">,</span>
       <span class="token property">"field"</span><span class="token operator">:</span>      <span class="token number">1</span><span class="token punctuation">,</span>
       <span class="token property">"searchtype"</span><span class="token operator">:</span> 'equals'<span class="token punctuation">,</span>
       <span class="token property">"value"</span><span class="token operator">:</span>      <span class="token number">1</span>
    <span class="token punctuation">}</span><span class="token punctuation">,</span> <span class="token punctuation">{</span>
       <span class="token property">"link"</span><span class="token operator">:</span>       'AND'<span class="token punctuation">,</span>
       <span class="token property">"criteria"</span> <span class="token operator">:</span> <span class="token punctuation">[</span>
          <span class="token punctuation">{</span>
             <span class="token property">"field"</span><span class="token operator">:</span>      <span class="token number">34</span><span class="token punctuation">,</span>
             <span class="token property">"searchtype"</span><span class="token operator">:</span> 'equals'<span class="token punctuation">,</span>
             <span class="token property">"value"</span><span class="token operator">:</span>      <span class="token number">1</span>
          <span class="token punctuation">}</span><span class="token punctuation">,</span> <span class="token punctuation">{</span>
             <span class="token property">"link"</span><span class="token operator">:</span>       'OR'<span class="token punctuation">,</span>
             <span class="token property">"field"</span><span class="token operator">:</span>      <span class="token number">35</span><span class="token punctuation">,</span>
             <span class="token property">"searchtype"</span><span class="token operator">:</span> 'equals'<span class="token punctuation">,</span>
             <span class="token property">"value"</span><span class="token operator">:</span>      <span class="token number">1</span>
          <span class="token punctuation">}</span>
       <span class="token punctuation">]</span>
    <span class="token punctuation">}</span>
 <span class="token punctuation">]</span>
...
</code></pre></li>
<li><p><em>metacriteria</em> (optional): array of meta-criterion objects to filter search. Optional.
                         A meta search is a link with another itemtype (ex: Computer with software).
<strong>Deprecated: Now criteria support meta flag, you should use it instead direct metacriteria option.</strong></p>

<p>Each meta-criterion object must provide:</p>

<ul>
<li><em>link</em>: logical operator in [AND, OR, AND NOT, AND NOT]. Mandatory.</li>
<li><em>itemtype</em>: second itemtype to link.</li>
<li><em>field</em>: id of the searchoption.</li>
<li><em>searchtype</em>: type of search in [contains¹, equals², notequals², lessthan, morethan, under, notunder].</li>
<li><em>value</em>: the value to search.</li>
</ul>

<p>Ex:</p>

<pre class="language-json" tabindex="0"><code class="language-json">...
<span class="token property">"metacriteria"</span><span class="token operator">:</span>
 <span class="token punctuation">[</span>
    <span class="token punctuation">{</span>
       <span class="token property">"link"</span><span class="token operator">:</span>       'AND'<span class="token punctuation">,</span>
       <span class="token property">"itemtype"</span><span class="token operator">:</span>   'Monitor'<span class="token punctuation">,</span>
       <span class="token property">"field"</span><span class="token operator">:</span>      <span class="token number">2</span><span class="token punctuation">,</span>
       <span class="token property">"searchtype"</span><span class="token operator">:</span> 'contains'<span class="token punctuation">,</span>
       <span class="token property">"value"</span><span class="token operator">:</span>      ''
    <span class="token punctuation">}</span><span class="token punctuation">,</span> <span class="token punctuation">{</span>
       <span class="token property">"link"</span><span class="token operator">:</span>       'AND'<span class="token punctuation">,</span>
       <span class="token property">"itemtype"</span><span class="token operator">:</span>   'Monitor'<span class="token punctuation">,</span>
       <span class="token property">"field"</span><span class="token operator">:</span>      <span class="token number">3</span><span class="token punctuation">,</span>
       <span class="token property">"searchtype"</span><span class="token operator">:</span> 'contains'<span class="token punctuation">,</span>
       <span class="token property">"value"</span><span class="token operator">:</span>      ''
     <span class="token punctuation">}</span>
 <span class="token punctuation">]</span>
...
</code></pre></li>
<li><p><em>sort</em> (default 1): id of the searchoption to sort by. Optional.</p></li>
<li><em>order</em> (default ASC): ASC - Ascending sort / DESC Descending sort. Optional.</li>
<li><em>range</em> (default 0-49): a string with a couple of number for start and end of pagination separated by a '-'. Ex: 150-199.
                     Optional.</li>
<li><em>forcedisplay</em>: array of columns to display (default empty = use display preferences and searched criteria).
             Some columns will be always presents (1: id, 2: name, 80: Entity).
             Optional.</li>
<li><em>rawdata</em> (default false): a boolean for displaying raws data of the Search engine of GLPI (like SQL request, full searchoptions, etc)</li>
<li><em>withindexes</em> (default false): a boolean to retrieve rows indexed by items id.
By default this option is set to false, because order of JSON objects (which are identified by index) cannot be garrantued  (from <a href="http://json.org/">http://json.org/</a> : An object is an unordered set of name/value pairs).
So, we provide arrays to guarantying sorted rows.</li>
<li><em>uid_cols</em> (default false): a boolean to identify cols by the 'uniqid' of the searchoptions instead of a numeric value (see <a href="#list-searchoptions">List searchOptions</a> and 'uid' field)</li>
<li><p><em>giveItems</em> (default false): a boolean to retrieve the data with the html parsed from core, new data are provided in data_html key.</p></li>
<li><p>¹ - <em>contains</em> will use a wildcard search per default. You can restrict at the beginning using the <em>^</em> character, and/or at the end using the <em>$</em> character.</p></li>
<li>² - <em>equals</em> and <em>notequals</em> are designed to be used with dropdowns. Do not expect those operators to search for a strictly equal value (see ¹ above).</li>
</ul></li>
<li><p><strong>Returns</strong>:</p>

<ul>
<li><p>200 (OK) with all rows data with this format:</p>

<pre class="language-json" tabindex="0"><code class="language-json"><span class="token punctuation">{</span>
    <span class="token property">"totalcount"</span><span class="token operator">:</span> <span class="token string">":numberofresults_without_pagination"</span><span class="token punctuation">,</span>
    <span class="token property">"range"</span><span class="token operator">:</span> <span class="token string">":start-:end"</span><span class="token punctuation">,</span>
    <span class="token property">"data"</span><span class="token operator">:</span> <span class="token punctuation">[</span>
        <span class="token punctuation">{</span>
            <span class="token property">":searchoptions_id"</span><span class="token operator">:</span> <span class="token string">"value"</span><span class="token punctuation">,</span>
            ...
        <span class="token punctuation">}</span><span class="token punctuation">,</span>
        <span class="token punctuation">{</span>
         ...
        <span class="token punctuation">}</span>
    <span class="token punctuation">]</span><span class="token punctuation">,</span>
    <span class="token property">"rawdata"</span><span class="token operator">:</span> <span class="token punctuation">{</span>
      ...
    <span class="token punctuation">}</span>
<span class="token punctuation">}</span>
</code></pre></li>
<li><p>206 (PARTIAL CONTENT) with rows data (pagination doesn't permit to display all rows).</p></li>
<li><p>401 (UNAUTHORIZED).</p>

<p>and theses headers:</p>

<ul>
<li><em>Content-Range</em> offset – limit / count</li>
<li><em>Accept-Range</em> itemtype max</li>
</ul></li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash"><span class="token function">curl</span> -g -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/search/Monitor?\
criteria\[0\]\[link\]\=AND\
\&amp;criteria\[0\]\[itemtype\]\=Monitor\
\&amp;criteria\[0\]\[field\]\=23\
\&amp;criteria\[0\]\[searchtype\]\=contains\
\&amp;criteria\[0\]\[value\]\=GSM\
\&amp;criteria\[1\]\[link\]\=AND\
\&amp;criteria\[1\]\[itemtype\]\=Monitor\
\&amp;criteria\[1\]\[field\]\=1\
\&amp;criteria\[1\]\[searchtype\]\=contains\
\&amp;criteria\[1\]\[value\]\=W2\
\&amp;range\=0-2\&amp;&amp;forcedisplay\[0\]\=1'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token operator">&lt;</span> Content-Range: <span class="token number">0</span>-2/2
<span class="token operator">&lt;</span> Accept-Range: <span class="token number">990</span>
<span class="token operator">&lt;</span> <span class="token punctuation">{</span><span class="token string">"totalcount"</span>:2,<span class="token string">"count"</span>:2,<span class="token string">"data"</span>:<span class="token punctuation">{</span><span class="token string">"11"</span>:<span class="token punctuation">{</span><span class="token string">"1"</span><span class="token builtin class-name">:</span><span class="token string">"W2242"</span>,<span class="token string">"80"</span><span class="token builtin class-name">:</span><span class="token string">"Root Entity"</span>,<span class="token string">"23"</span><span class="token builtin class-name">:</span><span class="token string">"GSM"</span><span class="token punctuation">}</span>,<span class="token string">"7"</span>:<span class="token punctuation">{</span><span class="token string">"1"</span><span class="token builtin class-name">:</span><span class="token string">"W2252"</span>,<span class="token string">"80"</span><span class="token builtin class-name">:</span><span class="token string">"Root Entity"</span>,<span class="token string">"23"</span><span class="token builtin class-name">:</span><span class="token string">"GSM"</span><span class="token punctuation">}</span><span class="token punctuation">}</span><span class="token punctuation">}</span>%
</code></pre>

<h2 id="add-items">Add item(s)</h2>

<ul>
<li><strong>URL</strong>: apirest.php/:itemtype/</li>
<li><strong>Description</strong>: Add an object (or multiple objects) into GLPI.</li>
<li><strong>Method</strong>: POST</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><p><strong>Parameters</strong>: (JSON Payload)</p>

<ul>
<li><em>input</em>: an object with fields of itemtype to be inserted.
      You can add several items in one action by passing an array of objects.
      Mandatory.</li>
</ul>

<p><strong>Important:</strong>
  In case of 'multipart/data' content_type (aka file upload), you should insert your parameters into
  a 'uploadManifest' parameter.
  Theses serialized data should be a JSON string.</p></li>
<li><p><strong>Returns</strong>:</p>

<ul>
<li>201 (OK) with id of added items.</li>
<li>207 (Multi-Status) with id of added items and errors.</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
<li>401 (UNAUTHORIZED).</li>
<li>And additional header can be provided on creation success:</li>
<li>Location when adding a single item.</li>
<li>Link on bulk addition.</li>
</ul></li>
</ul>

<p>Examples usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X POST <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"input": {"name": "My single computer", "serial": "12345"}}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Computer/'</span>

<span class="token operator">&lt;</span> <span class="token number">201</span> OK
<span class="token operator">&lt;</span> Location: http://path/to/glpi/api/Computer/15
<span class="token operator">&lt;</span> <span class="token punctuation">{</span><span class="token string">"id"</span><span class="token builtin class-name">:</span> <span class="token number">15</span><span class="token punctuation">}</span>


$ <span class="token function">curl</span> -X POST <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"input": [{"name": "My first computer", "serial": "12345"}, {"name": "My 2nd computer", "serial": "67890"}, {"name": "My 3rd computer", "serial": "qsd12sd"}]}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Computer/'</span>

<span class="token operator">&lt;</span> <span class="token number">207</span> OK
<span class="token operator">&lt;</span> Link: http://path/to/glpi/api/Computer/8,http://path/to/glpi/api/Computer/9
<span class="token operator">&lt;</span> <span class="token punctuation">[</span> <span class="token punctuation">{</span><span class="token string">"id"</span>:8, <span class="token string">"message"</span><span class="token builtin class-name">:</span> <span class="token string">""</span><span class="token punctuation">}</span>, <span class="token punctuation">{</span><span class="token string">"id"</span>:false, <span class="token string">"message"</span><span class="token builtin class-name">:</span> <span class="token string">"You don't have permission to perform this action."</span><span class="token punctuation">}</span>, <span class="token punctuation">{</span><span class="token string">"id"</span>:9, <span class="token string">"message"</span><span class="token builtin class-name">:</span> <span class="token string">""</span><span class="token punctuation">}</span> <span class="token punctuation">]</span>

</code></pre>

<p>Note: To upload a document see <a href="#upload-a-document-file">Upload a document file</a>.</p>

<h2 id="update-items">Update item(s)</h2>

<ul>
<li><strong>URL</strong>: apirest.php/:itemtype/:id</li>
<li><strong>Description</strong>: Update an object (or multiple objects) existing in GLPI.</li>
<li><strong>Method</strong>: PUT or PATCH</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (JSON Payload)

<ul>
<li><em>id</em>: the unique identifier of the itemtype passed in URL. You <strong>could skip</strong> this parameter by passing it in the input payload.</li>
<li><em>input</em>: Array of objects with fields of itemtype to be updated.
       Mandatory.
       You <strong>could provide</strong> in each object a key named 'id' to identify the item(s) to update.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) with update status for each item.</li>
<li>207 (Multi-Status) with id of added items and errors.</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
<li>401 (UNAUTHORIZED).</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X PUT <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"input": {"otherserial": "xcvbn"}}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Computer/10'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token punctuation">[</span><span class="token punctuation">{</span><span class="token string">"10"</span>:true, <span class="token string">"message"</span><span class="token builtin class-name">:</span> <span class="token string">""</span><span class="token punctuation">}</span><span class="token punctuation">]</span>


$ <span class="token function">curl</span> -X PUT <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"input": {"id": 11,  "otherserial": "abcde"}}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Computer/'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token punctuation">[</span><span class="token punctuation">{</span><span class="token string">"11"</span>:true, <span class="token string">"message"</span><span class="token builtin class-name">:</span> <span class="token string">""</span><span class="token punctuation">}</span><span class="token punctuation">]</span>


$ <span class="token function">curl</span> -X PUT <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"input": [{"id": 16,  "otherserial": "abcde"}, {"id": 17,  "otherserial": "fghij"}]}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Computer/'</span>

<span class="token operator">&lt;</span> <span class="token number">207</span> OK
<span class="token punctuation">[</span><span class="token punctuation">{</span><span class="token string">"8"</span>:true, <span class="token string">"message"</span><span class="token builtin class-name">:</span> <span class="token string">""</span><span class="token punctuation">}</span>,<span class="token punctuation">{</span><span class="token string">"2"</span>:false, <span class="token string">"message"</span><span class="token builtin class-name">:</span> <span class="token string">"Item not found"</span><span class="token punctuation">}</span><span class="token punctuation">]</span>
</code></pre>

<h2 id="delete-items">Delete item(s)</h2>

<ul>
<li><strong>URL</strong>: apirest.php/:itemtype/:id</li>
<li><strong>Description</strong>: Delete an object existing in GLPI.</li>
<li><strong>Method</strong>: DELETE</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><p><strong>Parameters</strong>: (query string)</p>

<ul>
<li><em>id</em>: unique identifier of the itemtype passed in the URL. You <strong>could skip</strong> this parameter by passing it in the input payload.
OR</li>
<li><em>input</em> Array of id who need to be deleted. This parameter is passed by payload.</li>
</ul>

<p>id parameter has precedence over input payload.</p>

<ul>
<li><em>force_purge</em> (default false): boolean, if the itemtype have a trashbin, you can force purge (delete finally).
             Optional.</li>
<li><em>history</em> (default true): boolean, set to false to disable saving of deletion in global history.
         Optional.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) <em>in case of multiple deletion</em>.</li>
<li>204 (No Content) <em>in case of single deletion</em>.</li>
<li>207 (Multi-Status) with id of deleted items and errors.</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
<li>401 (UNAUTHORIZED).</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X DELETE <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Computer/16?force_purge=true'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token punctuation">[</span><span class="token punctuation">{</span><span class="token string">"16"</span>:true, <span class="token string">"message"</span><span class="token builtin class-name">:</span> <span class="token string">""</span><span class="token punctuation">}</span><span class="token punctuation">]</span>

$ <span class="token function">curl</span> -X DELETE <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"input": {"id": 11}, "force_purge": true}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Computer/'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token punctuation">[</span><span class="token punctuation">{</span><span class="token string">"11"</span>:true, <span class="token string">"message"</span><span class="token builtin class-name">:</span> <span class="token string">""</span><span class="token punctuation">}</span><span class="token punctuation">]</span>


$ <span class="token function">curl</span> -X DELETE <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"input": [{"id": 16}, {"id": 17}]}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Computer/'</span>

<span class="token operator">&lt;</span> <span class="token number">207</span> OK
<span class="token punctuation">[</span><span class="token punctuation">{</span><span class="token string">"16"</span>:true, <span class="token string">"message"</span><span class="token builtin class-name">:</span> <span class="token string">""</span><span class="token punctuation">}</span>,<span class="token punctuation">{</span><span class="token string">"17"</span>:false, <span class="token string">"message"</span><span class="token builtin class-name">:</span> <span class="token string">"Item not found"</span><span class="token punctuation">}</span><span class="token punctuation">]</span>
</code></pre>

<h2 id="get-available-massive-actions-for-an-itemtype">Get available massive actions for an itemtype</h2>

<ul>
<li><strong>URL</strong>: apirest.php/getMassiveActions/:itemtype/</li>
<li><strong>Description</strong>: Show the availables massive actions for a given itemtype.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (query string)

<ul>
<li><em>is_deleted</em> (default false): Show specific actions for items in the trashbin.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK).</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
<li>401 (UNAUTHORIZED).</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getMassiveActions/Computer'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token punctuation">[</span>
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:update"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Update"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:clone"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Clone"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Infocom:activate"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Enable the financial and administrative information"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:delete"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Put in trashbin"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:add_transfer_list"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Add to transfer list"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Appliance:add_item"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Associate to an appliance"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Item_OperatingSystem:update"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Operating systems"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Computer_Item:add"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Connect"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Item_SoftwareVersion:add"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Install"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"KnowbaseItem_Item:add"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Link knowledgebase article"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Document_Item:add"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Add a document"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Document_Item:remove"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Remove a document"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Contract_Item:add"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Add a contract"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Contract_Item:remove"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Remove a contract"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:amend_comment"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Amend comment"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:add_note"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Add note"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Lock:unlock"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Unlock components"</span>
  <span class="token punctuation">}</span>
<span class="token punctuation">]</span>

$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getMassiveActions/Computer?is_deleted=1'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token punctuation">[</span>
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:purge_item_but_devices"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Delete permanently but keep devices"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:purge"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Delete permanently and remove devices"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:restore"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Restore"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Lock:unlock"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Unlock components"</span>
  <span class="token punctuation">}</span>
<span class="token punctuation">]</span>
</code></pre>

<h2 id="get-available-massive-actions-for-an-item">Get available massive actions for an item</h2>

<ul>
<li><strong>URL</strong>: apirest.php/getMassiveActions/:itemtype/:id</li>
<li><strong>Description</strong>: Show the availables massive actions for a given item.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (query string)</li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK).</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
<li>401 (UNAUTHORIZED).</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getMassiveActions/Computer/3'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token punctuation">[</span>
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:update"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Update"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:clone"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Clone"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Infocom:activate"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Enable the financial and administrative information"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:delete"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Put in trashbin"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:add_transfer_list"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Add to transfer list"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Appliance:add_item"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Associate to an appliance"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Item_OperatingSystem:update"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Operating systems"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Computer_Item:add"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Connect"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Item_SoftwareVersion:add"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Install"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"KnowbaseItem_Item:add"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Link knowledgebase article"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Document_Item:add"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Add a document"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Document_Item:remove"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Remove a document"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Contract_Item:add"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Add a contract"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Contract_Item:remove"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Remove a contract"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:amend_comment"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Amend comment"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:add_note"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Add note"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Lock:unlock"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Unlock components"</span>
  <span class="token punctuation">}</span>
<span class="token punctuation">]</span>


$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getMassiveActions/Computer/4'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token punctuation">[</span>
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:purge_item_but_devices"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Delete permanently but keep devices"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:purge"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Delete permanently and remove devices"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"MassiveAction:restore"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Restore"</span>
  <span class="token punctuation">}</span>,
  <span class="token punctuation">{</span>
    <span class="token string">"key"</span><span class="token builtin class-name">:</span> <span class="token string">"Lock:unlock"</span>,
    <span class="token string">"label"</span><span class="token builtin class-name">:</span> <span class="token string">"Unlock components"</span>
  <span class="token punctuation">}</span>
<span class="token punctuation">]</span>
</code></pre>

<h2 id="get-massive-action-parameters">Get massive action parameters</h2>

<ul>
<li><strong>URL</strong>: apirest.php/getMassiveActionParameters/:itemtype/</li>
<li><strong>Description</strong>: Show the availables parameters for a given massive action.

<ul>
<li>Warning: experimental endpoint, some required parameters may be missing from the returned content.</li>
</ul></li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (query string)</li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK).</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
<li>401 (UNAUTHORIZED).</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getMassiveActionParameters/Computer/MassiveAction:add_note'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token punctuation">[</span>
  <span class="token punctuation">{</span>
    <span class="token string">"name"</span><span class="token builtin class-name">:</span> <span class="token string">"add_note"</span>,
    <span class="token string">"type"</span><span class="token builtin class-name">:</span> <span class="token string">"text"</span>
  <span class="token punctuation">}</span>
<span class="token punctuation">]</span>

$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/getMassiveActionParameters/Computer/Contract_Item:add'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token punctuation">[</span>
  <span class="token punctuation">{</span>
    <span class="token string">"name"</span><span class="token builtin class-name">:</span> <span class="token string">"peer_contracts_id"</span>,
    <span class="token string">"type"</span><span class="token builtin class-name">:</span> <span class="token string">"dropdown"</span>
  <span class="token punctuation">}</span>
<span class="token punctuation">]</span>
</code></pre>

<h2 id="apply-massive-action">Apply massive action</h2>

<ul>
<li><strong>URL</strong>: apirest.php/applyMassiveAction/:itemtype/</li>
<li><strong>Description</strong>: Run the given massive action</li>
<li><strong>Method</strong>: POST</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Parameters</strong>: (json payload)

<ul>
<li><em>ids</em> items to execute the action on. Mandatory.</li>
<li><em>specific action parameters</em> some actions require specific parameters to run. You can check them through the 'getMassiveActionParameters' endpoint.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) All items were processed.</li>
<li>207 (Multi-Status) Not all items were successfully processed.</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
<li>401 (UNAUTHORIZED).</li>
<li>422 (Unprocessable Entity) All items failed to be processed.</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X POST <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"ids": [2, 3], "input": {"amendment": "newtext"}}'</span>
<span class="token string">'http://path/to/glpi/apirest.php/applyMassiveAction/Computer/MassiveAction:amend_comment'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
<span class="token punctuation">{</span>
  <span class="token string">"ok"</span><span class="token builtin class-name">:</span> <span class="token number">2</span>,
  <span class="token string">"ko"</span><span class="token builtin class-name">:</span> <span class="token number">0</span>,
  <span class="token string">"noright"</span><span class="token builtin class-name">:</span> <span class="token number">0</span>,
  <span class="token string">"messages"</span><span class="token builtin class-name">:</span> <span class="token punctuation">[</span><span class="token punctuation">]</span>
<span class="token punctuation">}</span>
</code></pre>

<h2 id="special-cases">Special cases</h2>

<h3 id="upload-a-document-file">Upload a document file</h3>

<p>See <a href="#add-items">Add item(s)</a> and apply specific instructions below.</p>

<p>Uploading a file requires use of 'multipart/data' content_type. The input data must be send in a 'uploadManifest' parameter and use the JSON format.</p>

<p>Examples usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X POST <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: multipart/form-data'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-F <span class="token string">'uploadManifest={"input": {"name": "Uploaded document", "_filename" : ["file.txt"]}};type=application/json'</span> <span class="token punctuation">\</span>
-F <span class="token string">'filename[0]=@file.txt'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Document/'</span>

<span class="token operator">&lt;</span> <span class="token number">201</span> OK
<span class="token operator">&lt;</span> Location: http://path/to/glpi/api/Document/1
<span class="token operator">&lt;</span> <span class="token punctuation">{</span><span class="token string">"id"</span><span class="token builtin class-name">:</span> <span class="token number">1</span>, <span class="token string">"message"</span><span class="token builtin class-name">:</span> <span class="token string">"Document move succeeded."</span>, <span class="token string">"upload_result"</span><span class="token builtin class-name">:</span> <span class="token punctuation">{</span><span class="token punctuation">..</span>.<span class="token punctuation">}</span><span class="token punctuation">}</span>

</code></pre>

<h3 id="download-a-document-file">Download a document file</h3>

<ul>
<li><strong>URL</strong>: apirest.php/Document/:id</li>
<li><strong>Description</strong>: Download a document.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
<li><em>Accept</em>: must be <strong>application/octet-stream</strong>. This header OR the parameter <em>alt</em> is mandatory</li>
</ul></li>
<li><p><strong>Parameters</strong>: (query string)</p>

<ul>
<li><em>id</em>: unique identifier of the itemtype passed in the URL. You <strong>could skip</strong> this parameter by passing it in the input payload.</li>
<li><em>alt</em>: must be 'media'. This parameter or the header <strong>Accept</strong> is mandatory.</li>
</ul>

<p>id parameter has precedence over input payload.</p></li>
<li><p><strong>Returns</strong>:</p>

<ul>
<li>200 (OK) <em>in case of multiple deletion</em>.</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
<li>401 (UNAUTHORIZED).</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-H <span class="token string">"Accept: application/octet-stream"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"input": {"id": 11}}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Document/'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
</code></pre>

<p>The body of the answer contains the raw file attached to the document.</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
-d <span class="token string">'{"input": {"id": 11}}'</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/Document/&amp;alt=media'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
</code></pre>

<p>The body of the answer contains the raw file attached to the document.</p>

<h3 id="get-a-user%27s-profile-picture">Get a user's profile picture</h3>

<ul>
<li><strong>URL</strong>: apirest.php/User/:id/Picture</li>
<li><strong>Description</strong>: Get a user's profile picture.</li>
<li><strong>Method</strong>: GET</li>
<li><strong>Parameters</strong>: (Headers)

<ul>
<li><em>Session-Token</em>: session var provided by <a href="#init-session">initSession</a> endpoint. Mandatory.</li>
<li><em>App-Token</em>: authorization string provided by the GLPI API configuration. Optional.</li>
</ul></li>
<li><strong>Returns</strong>:

<ul>
<li>200 (OK) with the raw image in the request body.</li>
<li>204 (No content) if the request is correct but the specified user doesn't have a profile picture.</li>
<li>400 (Bad Request) with a message indicating an error in input parameter.</li>
</ul></li>
</ul>

<p>Example usage (CURL):</p>

<pre class="language-bash" tabindex="0"><code class="language-bash">$ <span class="token function">curl</span> -X GET <span class="token punctuation">\</span>
-H <span class="token string">'Content-Type: application/json'</span> <span class="token punctuation">\</span>
-H <span class="token string">"Session-Token: 83af7e620c83a50a18d3eac2f6ed05a3ca0bea62"</span> <span class="token punctuation">\</span>
-H <span class="token string">"App-Token: f7g3csp8mgatg5ebc5elnazakw20i9fyev1qopya7"</span> <span class="token punctuation">\</span>
<span class="token string">'http://path/to/glpi/apirest.php/User/2/Picture/'</span>

<span class="token operator">&lt;</span> <span class="token number">200</span> OK
</code></pre>

<p>The body of the answer contains the raw image.</p>

<h3 id="sanitized-content">Sanitized content</h3>

<p>By default, the API will return sanitized content.<br>
This mean that all HTML special characters will be encoded.<br>
You can disable this feature by adding the following header to your request:</p>

<pre><code>X-GLPI-Sanitized-Content: false
</code></pre>

<h2 id="errors">Errors</h2>

<h3 id="error_item_not_found">ERROR_ITEM_NOT_FOUND</h3>

<p>The desired resource (itemtype-id) was not found in the GLPI database.</p>

<h3 id="error_bad_array">ERROR_BAD_ARRAY</h3>

<p>The HTTP body must be an an array of objects.</p>

<h3 id="error_method_not_allowed">ERROR_METHOD_NOT_ALLOWED</h3>

<p>You specified an inexistent or not not allowed resource.</p>

<h3 id="error_right_missing">ERROR_RIGHT_MISSING</h3>

<p>The current logged user miss rights in his profile to do the provided action.
Alter this profile or choose a new one for the user in GLPI main interface.</p>

<h3 id="error_session_token_invalid">ERROR_SESSION_TOKEN_INVALID</h3>

<p>The Session-Token provided in header is invalid.
You should redo an <a href="#init-session">Init session</a> request.</p>

<h3 id="error_session_token_missing">ERROR_SESSION_TOKEN_MISSING</h3>

<p>You miss to provide Session-Token in header of your HTTP request.</p>

<h3 id="error_app_token_parameters_missing">ERROR_APP_TOKEN_PARAMETERS_MISSING</h3>

<p>The current API requires an App-Token header for using its methods.</p>

<h3 id="error_wrong_app_token_parameter">ERROR_WRONG_APP_TOKEN_PARAMETER</h3>

<p>It seems the provided application token doesn't exists in GLPI API configuration.</p>

<h3 id="error_not_deleted">ERROR_NOT_DELETED</h3>

<p>You must mark the item for deletion before actually deleting it</p>

<h3 id="error_not_allowed_ip">ERROR_NOT_ALLOWED_IP</h3>

<p>We can't find an active client defined in configuration for your IP.
Go to the GLPI Configuration &gt; Setup menu and API tab to check IP access.</p>

<h3 id="error_login_parameters_missing">ERROR_LOGIN_PARAMETERS_MISSING</h3>

<p>One of theses parameter(s) is missing:</p>

<ul>
<li>login and password</li>
<li>or user_token</li>
</ul>

<h3 id="error_login_with_credentials_disabled">ERROR_LOGIN_WITH_CREDENTIALS_DISABLED</h3>

<p>The GLPI setup forbid the login with credentials, you must login with your user_token instead.
See your personal preferences page or setup API access in GLPI main interface.</p>

<h3 id="error_glpi_login_user_token">ERROR_GLPI_LOGIN_USER_TOKEN</h3>

<p>The provided user_token seems invalid.
Check your personal preferences page in GLPI main interface.</p>

<h3 id="error_glpi_login">ERROR_GLPI_LOGIN</h3>

<p>We cannot login you into GLPI. This error is not relative to API but GLPI core.
Check the user administration and the GLPI logs files (in files/_logs directory).</p>

<h3 id="error_itemtype_not_found_nor_commondbtm">ERROR_ITEMTYPE_NOT_FOUND_NOR_COMMONDBTM</h3>

<p>You asked a inexistent resource (endpoint). It's not a predefined (initSession, getFullSession, etc) nor a GLPI CommonDBTM resources.</p>

<p>See this documentation for predefined ones or <a href="https://forge.glpi-project.org/apidoc/class-CommonDBTM.html">List itemtypes</a> for available resources</p>

<h3 id="error_sql">ERROR_SQL</h3>

<p>We suspect an SQL error.
This error is not relative to API but to GLPI core.
Check the GLPI logs files (in files/_logs directory).</p>

<h3 id="error_range_exceed_total">ERROR_RANGE_EXCEED_TOTAL</h3>

<p>The range parameter you provided is superior to the total count of available data.</p>

<h3 id="error_glpi_add">ERROR_GLPI_ADD</h3>

<p>We cannot add the object to GLPI. This error is not relative to API but to GLPI core.
Check the GLPI logs files (in files/_logs directory).</p>

<h3 id="error_glpi_partial_add">ERROR_GLPI_PARTIAL_ADD</h3>

<p>Some of the object you wanted to add triggers an error.
Maybe a missing field or rights.
You'll find with this error a collection of results.</p>

<h3 id="error_glpi_update">ERROR_GLPI_UPDATE</h3>

<p>We cannot update the object to GLPI. This error is not relative to API but to GLPI core.
Check the GLPI logs files (in files/_logs directory).</p>

<h3 id="error_glpi_partial_update">ERROR_GLPI_PARTIAL_UPDATE</h3>

<p>Some of the object you wanted to update triggers an error.
Maybe a missing field or rights.
You'll find with this error a collection of results.</p>

<h3 id="error_glpi_delete">ERROR_GLPI_DELETE</h3>

<p>We cannot delete the object to GLPI. This error is not relative to API but to GLPI core.
Check the GLPI logs files (in files/_logs directory).</p>

<h3 id="error_glpi_partial_delete">ERROR_GLPI_PARTIAL_DELETE</h3>

<p>Some of the objects you want to delete triggers an error, maybe a missing field or rights.
You'll find with this error, a collection of results.</p>

<h3 id="error_massiveaction_key">ERROR_MASSIVEACTION_KEY</h3>

<p>Missing or invalid massive action key.
Run 'getMassiveActions' endpoint to see available keys.</p>

<h3 id="error_massiveaction_no_ids">ERROR_MASSIVEACTION_NO_IDS</h3>

<p>No ids supplied when trying to run a massive action.</p>

<h2 id="servers-configuration">Servers configuration</h2>

<p>By default, you can use <a href="http://path/to/glpi/apirest.php">http://path/to/glpi/apirest.php</a> without any additional configuration.</p>

<p>You'll find below some examples to configure your web server to redirect your <a href="http://.../glpi/api/">http://.../glpi/api/</a> URL to the apirest.php file.</p>

<h3 id="apache-httpd">Apache Httpd</h3>

<p>We provide in root .htaccess of GLPI an example to enable API URL rewriting.</p>

<p>You need to uncomment (removing #) theses lines:</p>

<pre class="language-apacheconf" tabindex="0"><code class="language-apacheconf"><span class="token comment">#&lt;IfModule mod_rewrite.c&gt;</span>
<span class="token comment">#   RewriteEngine On</span>
<span class="token comment">#   RewriteCond %{REQUEST_FILENAME} !-f</span>
<span class="token comment">#   RewriteCond %{REQUEST_FILENAME} !-d</span>
<span class="token comment">#   RewriteRule api/(.*)$ apirest.php/$1</span>
<span class="token comment">#&lt;/IfModule&gt;</span>
</code></pre>

<p>By enabling URL rewriting, you could use API with this URL : <a href="http://path/to/glpi/api/">http://path/to/glpi/api/</a>.
You need also to enable rewrite module in apache httpd and permit GLPI's .htaccess to override server configuration (see AllowOverride directive).</p>

<p><strong>Note for apache+fpm users:</strong></p>

<p>You may have difficulties to pass Authorization header in this configuration.
You have two options :</p>

<ul>
<li>pass the <code>user_token</code> or credentials (login/password) in the http query (as GET parameters).</li>
<li>add env to your virtualhost: <code>SetEnvIf Authorization "(.*)" HTTP_AUTHORIZATION=$1</code>.</li>
</ul>

<h3 id="nginx">Nginx</h3>

<p>This example of configuration was achieved on ubuntu with php7 fpm.</p>

<pre class="language-nginx" tabindex="0"><code class="language-nginx"><span class="token directive"><span class="token keyword">server</span></span> <span class="token punctuation">{</span>
   <span class="token directive"><span class="token keyword">listen</span> <span class="token number">80</span> default_server</span><span class="token punctuation">;</span>
   <span class="token directive"><span class="token keyword">listen</span> [::]:80 default_server</span><span class="token punctuation">;</span>

   <span class="token comment"># change here to match your GLPI directory</span>
   <span class="token directive"><span class="token keyword">root</span> /var/www/html/glpi/</span><span class="token punctuation">;</span>

   <span class="token directive"><span class="token keyword">index</span> index.html index.htm index.nginx-debian.html index.php</span><span class="token punctuation">;</span>

   <span class="token directive"><span class="token keyword">server_name</span> localhost</span><span class="token punctuation">;</span>

   <span class="token directive"><span class="token keyword">location</span> /</span> <span class="token punctuation">{</span>
      <span class="token directive"><span class="token keyword">try_files</span> <span class="token variable">$uri</span> <span class="token variable">$uri</span>/ =404</span><span class="token punctuation">;</span>
      <span class="token directive"><span class="token keyword">autoindex</span> <span class="token boolean">on</span></span><span class="token punctuation">;</span>
   <span class="token punctuation">}</span>

   <span class="token directive"><span class="token keyword">location</span> /api</span> <span class="token punctuation">{</span>
      <span class="token directive"><span class="token keyword">rewrite</span> ^/api/(.*)$ /apirest.php/<span class="token variable">$1</span> last</span><span class="token punctuation">;</span>
   <span class="token punctuation">}</span>

   <span class="token directive"><span class="token keyword">location</span> ~ [^/]\.php(/|$)</span> <span class="token punctuation">{</span>
      <span class="token directive"><span class="token keyword">fastcgi_pass</span> unix:/run/php/php7.0-fpm.sock</span><span class="token punctuation">;</span>

      <span class="token comment"># regex to split $uri to $fastcgi_script_name and $fastcgi_path</span>
      <span class="token directive"><span class="token keyword">fastcgi_split_path_info</span> ^(.+\.php)(/.+)$</span><span class="token punctuation">;</span>

      <span class="token comment"># Check that the PHP script exists before passing it</span>
      <span class="token directive"><span class="token keyword">try_files</span> <span class="token variable">$fastcgi_script_name</span> =404</span><span class="token punctuation">;</span>

      <span class="token comment"># Bypass the fact that try_files resets $fastcgi_path_info</span>
      <span class="token comment"># # see: http://trac.nginx.org/nginx/ticket/321</span>
      <span class="token directive"><span class="token keyword">set</span> <span class="token variable">$path_info</span> <span class="token variable">$fastcgi_path_info</span></span><span class="token punctuation">;</span>
      <span class="token directive"><span class="token keyword">fastcgi_param</span>  PATH_INFO <span class="token variable">$path_info</span></span><span class="token punctuation">;</span>

      <span class="token directive"><span class="token keyword">fastcgi_param</span>  PATH_TRANSLATED    <span class="token variable">$document_root</span><span class="token variable">$fastcgi_script_name</span></span><span class="token punctuation">;</span>
      <span class="token directive"><span class="token keyword">fastcgi_param</span>  SCRIPT_FILENAME    <span class="token variable">$document_root</span><span class="token variable">$fastcgi_script_name</span></span><span class="token punctuation">;</span>

      <span class="token directive"><span class="token keyword">include</span> fastcgi_params</span><span class="token punctuation">;</span>

      <span class="token comment"># allow directory index</span>
      <span class="token directive"><span class="token keyword">fastcgi_index</span> index.php</span><span class="token punctuation">;</span>
   <span class="token punctuation">}</span>
<span class="token punctuation">}</span>

</code></pre>
</div><script type="text/javascript">
//<![CDATA[


         var CFG_GLPI  = {"languages":{"ar_SA":["\u0627\u0644\u0639\u064e\u0631\u064e\u0628\u0650\u064a\u064e\u0651\u0629\u064f","ar_SA.mo","ar","ar","arabic",103],"bg_BG":["\u0411\u044a\u043b\u0433\u0430\u0440\u0441\u043a\u0438","bg_BG.mo","bg","bg","bulgarian",2],"id_ID":["Bahasa Indonesia","id_ID.mo","id","id","indonesian",2],"ms_MY":["Bahasa Melayu","ms_MY.mo","ms","ms","malay",2],"ca_ES":["Catal\u00e0","ca_ES.mo","ca","ca","catalan",2],"cs_CZ":["\u010ce\u0161tina","cs_CZ.mo","cs","cs","czech",10],"de_DE":["Deutsch","de_DE.mo","de","de","german",2],"da_DK":["Dansk","da_DK.mo","da","da","danish",2],"et_EE":["Eesti","et_EE.mo","et","et","estonian",2],"en_GB":["English","en_GB.mo","en-GB","en","english",2],"en_US":["English (US)","en_US.mo","en-GB","en","english",2],"es_AR":["Espa\u00f1ol (Argentina)","es_AR.mo","es","es","spanish",2],"es_EC":["Espa\u00f1ol (Ecuador)","es_EC.mo","es","es","spanish",2],"es_CO":["Espa\u00f1ol (Colombia)","es_CO.mo","es","es","spanish",2],"es_ES":["Espa\u00f1ol (Espa\u00f1a)","es_ES.mo","es","es","spanish",2],"es_419":["Espa\u00f1ol (Am\u00e9rica Latina)","es_419.mo","es","es","spanish",2],"es_MX":["Espa\u00f1ol (Mexico)","es_MX.mo","es","es","spanish",2],"es_VE":["Espa\u00f1ol (Venezuela)","es_VE.mo","es","es","spanish",2],"eu_ES":["Euskara","eu_ES.mo","eu","eu","basque",2],"fr_FR":["Fran\u00e7ais","fr_FR.mo","fr","fr","french",2],"fr_CA":["Fran\u00e7ais (Canada)","fr_CA.mo","fr","fr","french",2],"fr_BE":["Fran\u00e7ais (Belgique)","fr_BE.mo","fr","fr","french",2],"gl_ES":["Galego","gl_ES.mo","gl","gl","galician",2],"el_GR":["\u0395\u03bb\u03bb\u03b7\u03bd\u03b9\u03ba\u03ac","el_GR.mo","el","el","greek",2],"he_IL":["\u05e2\u05d1\u05e8\u05d9\u05ea","he_IL.mo","he","he","hebrew",2],"hi_IN":["\u0939\u093f\u0928\u094d\u0926\u0940","hi_IN.mo","hi","hi_IN","hindi",2],"hr_HR":["Hrvatski","hr_HR.mo","hr","hr","croatian",2],"hu_HU":["Magyar","hu_HU.mo","hu","hu","hungarian",2],"it_IT":["Italiano","it_IT.mo","it","it","italian",2],"kn":["\u0c95\u0ca8\u0ccd\u0ca8\u0ca1","kn.mo","en-GB","en","kannada",2],"lv_LV":["Latvie\u0161u","lv_LV.mo","lv","lv","latvian",2],"lt_LT":["Lietuvi\u0173","lt_LT.mo","lt","lt","lithuanian",2],"mn_MN":["\u041c\u043e\u043d\u0433\u043e\u043b \u0445\u044d\u043b","mn_MN.mo","mn","mn","mongolian",2],"nl_NL":["Nederlands","nl_NL.mo","nl","nl","dutch",2],"nl_BE":["Flemish","nl_BE.mo","nl","nl","flemish",2],"nb_NO":["Norsk (Bokm\u00e5l)","nb_NO.mo","no","nb","norwegian",2],"nn_NO":["Norsk (Nynorsk)","nn_NO.mo","no","nn","norwegian",2],"fa_IR":["\u0641\u0627\u0631\u0633\u06cc","fa_IR.mo","fa","fa","persian",2],"pl_PL":["Polski","pl_PL.mo","pl","pl","polish",2],"pt_PT":["Portugu\u00eas","pt_PT.mo","pt","pt","portuguese",2],"pt_BR":["Portugu\u00eas do Brasil","pt_BR.mo","pt-BR","pt","brazilian portuguese",2],"ro_RO":["Rom\u00e2n\u0103","ro_RO.mo","ro","en","romanian",2],"ru_RU":["\u0420\u0443\u0441\u0441\u043a\u0438\u0439","ru_RU.mo","ru","ru","russian",2],"sk_SK":["Sloven\u010dina","sk_SK.mo","sk","sk","slovak",10],"sl_SI":["Sloven\u0161\u010dina","sl_SI.mo","sl","sl","slovenian slovene",2],"sr_RS":["Srpski","sr_RS.mo","sr","sr","serbian",2],"fi_FI":["Suomi","fi_FI.mo","fi","fi","finish",2],"sv_SE":["Svenska","sv_SE.mo","sv","sv","swedish",2],"vi_VN":["Ti\u1ebfng Vi\u1ec7t","vi_VN.mo","vi","vi","vietnamese",2],"th_TH":["\u0e20\u0e32\u0e29\u0e32\u0e44\u0e17\u0e22","th_TH.mo","th","th","thai",2],"tr_TR":["T\u00fcrk\u00e7e","tr_TR.mo","tr","tr","turkish",2],"uk_UA":["\u0423\u043a\u0440\u0430\u0457\u043d\u0441\u044c\u043a\u0430","uk_UA.mo","uk","en","ukrainian",2],"ja_JP":["\u65e5\u672c\u8a9e","ja_JP.mo","ja","ja","japanese",2],"zh_CN":["\u7b80\u4f53\u4e2d\u6587","zh_CN.mo","zh-CN","zh","chinese",2],"zh_TW":["\u7e41\u9ad4\u4e2d\u6587","zh_TW.mo","zh-TW","zh","chinese",2],"ko_KR":["\ud55c\uad6d\/\u97d3\u570b","ko_KR.mo","ko","ko","korean",1],"zh_HK":["\u9999\u6e2f","zh_HK.mo","zh-HK","zh","chinese",2],"be_BY":["Belarussian","be_BY.mo","be","be","belarussian",3],"is_IS":["\u00edslenska","is_IS.mo","is","en","icelandic",2],"eo":["Esperanto","eo.mo","eo","en","esperanto",2],"es_CL":["Espa\u00f1ol chileno","es_CL","es","es","spanish chilean",2]},"glpitables":[],"unicity_types":["Budget","Computer","Contact","Contract","Infocom","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software","SoftwareLicense","Supplier","User","Certificate","Rack","Enclosure","PDU","Cluster","Item_DeviceSimcard"],"state_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer","SoftwareLicense","Certificate","Enclosure","PDU","Line","Rack","SoftwareVersion","Cluster","Contract","Appliance","DatabaseInstance","Cable"],"asset_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer","SoftwareLicense","Certificate"],"project_asset_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software","DeviceMotherboard","DeviceProcessor","DeviceMemory","DeviceHardDrive","DeviceNetworkCard","DeviceDrive","DeviceControl","DeviceGraphicCard","DeviceSoundCard","DevicePci","DeviceCase","DevicePowerSupply","DeviceGeneric","DeviceBattery","DeviceFirmware","DeviceCamera","Certificate"],"document_types":["Budget","CartridgeItem","Change","Computer","ConsumableItem","Contact","Contract","Document","Entity","KnowbaseItem","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Problem","Project","ProjectTask","Reminder","Software","Line","SoftwareLicense","Supplier","Ticket","User","Certificate","Cluster","ITILFollowup","ITILSolution","ChangeTask","ProblemTask","TicketTask","Appliance","DatabaseInstance","PluginFormcreatorFormAnswer","PluginFormcreatorForm"],"consumables_types":["Group","User"],"report_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Project","Software","SoftwareLicense","Certificate"],"directconnect_types":["Monitor","Peripheral","Phone","Printer"],"infocom_types":["Cartridge","CartridgeItem","Computer","Consumable","ConsumableItem","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software","SoftwareLicense","Line","Certificate","Domain","Appliance","Item_DeviceSimcard","Rack","Enclosure","PDU","PassiveDCEquipment","DatabaseInstance","Cable"],"reservation_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software"],"linkuser_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software","SoftwareLicense","Certificate","Appliance","Item_DeviceSimcard","Line"],"linkgroup_types":["Computer","Consumable","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software","SoftwareLicense","Certificate","Appliance","Item_DeviceSimcard","Line"],"linkuser_tech_types":["Computer","ConsumableItem","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software","SoftwareLicense","Certificate","Appliance","DatabaseInstance"],"linkgroup_tech_types":["Computer","ConsumableItem","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software","SoftwareLicense","Certificate","Appliance","DatabaseInstance"],"location_types":["Budget","CartridgeItem","ConsumableItem","Computer","Monitor","Glpi\\Socket","NetworkEquipment","Peripheral","Phone","Printer","Software","SoftwareLicense","Ticket","User","Certificate","Item_DeviceSimcard","Line","Appliance","PassiveDCEquipment","DataCenter","DCRoom","Rack","Enclosure","PDU","Item_DeviceMotherboard","Item_DeviceFirmware","Item_DeviceProcessor","Item_DeviceMemory","Item_DeviceHardDrive","Item_DeviceNetworkCard","Item_DeviceDrive","Item_DeviceBattery","Item_DeviceGraphicCard","Item_DeviceSoundCard","Item_DeviceControl","Item_DevicePci","Item_DeviceCase","Item_DevicePowerSupply","Item_DeviceGeneric","Item_DeviceSimcard","Item_DeviceSensor","Item_DeviceCamera"],"ticket_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software","SoftwareLicense","Certificate","Line","DCRoom","Rack","Enclosure","Cluster","PDU","Domain","DomainRecord","Appliance","Item_DeviceSimcard","PassiveDCEquipment","DatabaseInstance","Database","Cable","PluginFormcreatorFormAnswer"],"link_types":["Budget","CartridgeItem","Computer","ConsumableItem","Contact","Contract","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software","Supplier","User","Certificate","Cluster","DCRoom","Domain","Appliance","DatabaseInstance"],"dictionnary_types":["ComputerModel","ComputerType","Manufacturer","MonitorModel","MonitorType","NetworkEquipmentModel","NetworkEquipmentType","OperatingSystem","OperatingSystemServicePack","OperatingSystemVersion","PeripheralModel","PeripheralType","PhoneModel","PhoneType","Printer","PrinterModel","PrinterType","Software","OperatingSystemArchitecture","OperatingSystemKernel","OperatingSystemKernelVersion","OperatingSystemEdition","ImageResolution","ImageFormat","DatabaseInstanceType","Glpi\\SocketModel","CableType"],"helpdesk_visible_types":["Software","Appliance","Database"],"networkport_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Enclosure","PDU","Cluster","Unmanaged"],"networkport_instantiations":["NetworkPortEthernet","NetworkPortWifi","NetworkPortAggregate","NetworkPortAlias","NetworkPortDialup","NetworkPortLocal","NetworkPortFiberchannel"],"device_types":["DeviceMotherboard","DeviceFirmware","DeviceProcessor","DeviceMemory","DeviceHardDrive","DeviceNetworkCard","DeviceDrive","DeviceBattery","DeviceGraphicCard","DeviceSoundCard","DeviceControl","DevicePci","DeviceCase","DevicePowerSupply","DeviceGeneric","DeviceSimcard","DeviceSensor","DeviceCamera"],"socket_types":["Computer","NetworkEquipment","Peripheral","Phone","Printer","PassiveDCEquipment"],"itemdevices":["Item_DeviceMotherboard","Item_DeviceFirmware","Item_DeviceProcessor","Item_DeviceMemory","Item_DeviceHardDrive","Item_DeviceNetworkCard","Item_DeviceDrive","Item_DeviceBattery","Item_DeviceGraphicCard","Item_DeviceSoundCard","Item_DeviceControl","Item_DevicePci","Item_DeviceCase","Item_DevicePowerSupply","Item_DeviceGeneric","Item_DeviceSimcard","Item_DeviceSensor","Item_DeviceCamera"],"itemdevices_types":["Computer","NetworkEquipment","Peripheral","Phone","Printer","Enclosure"],"itemdevices_itemaffinity":["Computer"],"itemdevicememory_types":["Computer","NetworkEquipment","Peripheral","Printer","Phone"],"itemdevicepowersupply_types":["Computer","NetworkEquipment","Enclosure","Phone"],"itemdevicenetworkcard_types":["Computer","NetworkEquipment","Peripheral","Phone","Printer"],"itemdeviceharddrive_types":["Computer","Peripheral","NetworkEquipment","Printer","Phone"],"itemdevicebattery_types":["Computer","Peripheral","Phone","Printer"],"itemdevicefirmware_types":["Computer","Peripheral","Phone","NetworkEquipment","Printer"],"itemdevicesimcard_types":["Computer","Peripheral","Phone","NetworkEquipment","Printer"],"itemdevicegeneric_types":["*"],"itemdevicepci_types":["*"],"itemdevicesensor_types":["Computer","Peripheral","Phone"],"itemdeviceprocessor_types":["Computer","Phone"],"itemdevicesoundcard_types":["Computer"],"itemdevicegraphiccard_types":["Computer","Phone"],"itemdevicemotherboard_types":["Computer","Phone"],"itemdevicecamera_types":["Computer","Phone"],"notificationtemplates_types":["CartridgeItem","Change","ConsumableItem","Contract","CronTask","DBConnection","FieldUnicity","Infocom","MailCollector","ObjectLock","PlanningRecall","Problem","Project","ProjectTask","Reservation","SoftwareLicense","Ticket","User","SavedSearch_Alert","Certificate","Glpi\\Marketplace\\Controller","Domain","PluginFormcreatorFormAnswer"],"contract_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Project","Line","Software","SoftwareLicense","Certificate","DCRoom","Rack","Enclosure","Cluster","PDU","Appliance","Domain","DatabaseInstance","Item_DeviceMotherboard","Item_DeviceFirmware","Item_DeviceProcessor","Item_DeviceMemory","Item_DeviceHardDrive","Item_DeviceNetworkCard","Item_DeviceDrive","Item_DeviceBattery","Item_DeviceGraphicCard","Item_DeviceSoundCard","Item_DeviceControl","Item_DevicePci","Item_DeviceCase","Item_DevicePowerSupply","Item_DeviceGeneric","Item_DeviceSimcard","Item_DeviceSensor","Item_DeviceCamera"],"union_search_type":{"ReservationItem":"reservation_types","AllAssets":"asset_types"},"systeminformations_types":["AuthLDAP","DBConnection","MailCollector","Plugin"],"rulecollections_types":["RuleImportAssetCollection","RuleImportEntityCollection","RuleLocationCollection","RuleMailCollectorCollection","RuleRightCollection","RuleSoftwareCategoryCollection","RuleTicketCollection","RuleAssetCollection"],"planning_types":["ChangeTask","ProblemTask","Reminder","TicketTask","ProjectTask","PlanningExternalEvent"],"planning_add_types":["PlanningExternalEvent"],"caldav_supported_components":["VEVENT","VJOURNAL"],"globalsearch_types":["Computer","Contact","Contract","Document","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software","SoftwareLicense","Ticket","Problem","Change","User","Group","Project","Supplier","Budget","Certificate","Line","Datacenter","DCRoom","Enclosure","PDU","Rack","Cluster","PassiveDCEquipment","Domain","Appliance"],"number_format":"0","decimal_number":"2","debug_lang":1,"debug_vars":1,"debug_sql":1,"user_pref_field":["backcreated","csv_delimiter","date_format","default_requesttypes_id","display_count_on_home","duedatecritical_color","duedatecritical_less","duedatecritical_unit","duedateok_color","duedatewarning_color","duedatewarning_less","duedatewarning_unit","followup_private","is_ids_visible","keep_devices_when_purging_item","language","list_limit","lock_autolock_mode","lock_directunlock_notification","names_format","notification_to_myself","number_format","pdffont","priority_1","priority_2","priority_3","priority_4","priority_5","priority_6","refresh_views","set_default_tech","set_default_requester","show_count_on_tabs","show_jobs_at_login","task_private","task_state","use_flat_dropdowntree","palette","page_layout","highcontrast_css","default_dashboard_central","default_dashboard_assets","default_dashboard_helpdesk","default_dashboard_mini_ticket","default_central_tab","fold_menu","fold_search","savedsearches_pinned","richtext_layout","timeline_order","itil_layout"],"lock_lockable_objects":["Budget","Change","Contact","Contract","Document","CartridgeItem","Computer","ConsumableItem","Entity","Group","KnowbaseItem","Line","Link","Monitor","NetworkEquipment","NetworkName","Peripheral","Phone","Printer","Problem","Profile","Project","Reminder","RSSFeed","Software","Supplier","Ticket","User","SoftwareLicense","Certificate"],"inventory_types":["Computer","Phone","Printer","NetworkEquipment"],"inventory_lockable_objects":["Computer_Item","Item_SoftwareLicense","Item_SoftwareVersion","Item_Disk","ComputerVirtualMachine","ComputerAntivirus","NetworkPort","NetworkName","IPAddress","Item_OperatingSystem","Item_DeviceBattery","Item_DeviceCase","Item_DeviceControl","Item_DeviceDrive","Item_DeviceFirmware","Item_DeviceGeneric","Item_DeviceGraphicCard","Item_DeviceHardDrive","Item_DeviceMemory","Item_DeviceMotherboard","Item_DeviceNetworkCard","Item_DevicePci","Item_DevicePowerSupply","Item_DeviceProcessor","Item_DeviceSensor","Item_DeviceSimcard","Item_DeviceSoundCard","DatabaseInstance","Item_RemoteManagement","Monitor"],"kb_types":["Budget","Change","Computer","Contract","Entity","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Problem","Project","Software","SoftwareLicense","Supplier","Ticket","Certificate","Appliance","DatabaseInstance"],"certificate_types":["Computer","NetworkEquipment","Peripheral","Phone","Printer","SoftwareLicense","User","Domain","Appliance","DatabaseInstance"],"rackable_types":["Computer","Monitor","NetworkEquipment","Peripheral","Enclosure","PDU","PassiveDCEquipment"],"cluster_types":["Computer","NetworkEquipment"],"operatingsystem_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer"],"software_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer"],"kanban_types":["Project"],"domain_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software","Appliance","Certificate","DatabaseInstance","Database"],"appliance_types":["Computer","Monitor","NetworkEquipment","Peripheral","Phone","Printer","Software","Appliance","Cluster","DatabaseInstance","Database"],"appliance_relation_types":["Location","Network","Domain"],"remote_management_types":["Computer","Phone"],"databaseinstance_types":["Computer"],"agent_types":["Computer","Phone"],"javascript":{"central":{"central":["fullcalendar","planning","masonry","tinymce","dashboard","gridstack","charts","clipboard","sortable"]},"assets":{"dashboard":["dashboard","gridstack","charts","clipboard","sortable"],"rack":["gridstack","rack"],"printer":["dashboard","gridstack","charts","clipboard","sortable","fullcalendar","reservations"],"cable":["cable"],"socket":["cable"],"networkport":["dashboard","gridstack","charts","clipboard","sortable"],"computer":["fullcalendar","reservations"],"monitor":["fullcalendar","reservations"],"networkequipment":["fullcalendar","reservations"],"peripheral":["fullcalendar","reservations"],"phone":["fullcalendar","reservations"],"software":["fullcalendar","reservations"]},"helpdesk":{"dashboard":["dashboard","gridstack","charts","clipboard","sortable"],"planning":["clipboard","fullcalendar","tinymce","planning"],"ticket":["rateit","tinymce","kanban","dashboard","gridstack","charts","clipboard","sortable"],"problem":["tinymce","kanban","sortable"],"change":["tinymce","kanban","sortable"],"stat":["charts"],"pluginformcreatorformlist":["gridstack"],"pluginformcreatorissue":["photoswipe"]},"tools":{"project":["kanban","tinymce","sortable"],"knowbaseitem":["tinymce"],"knowbaseitemtranslation":["tinymce"],"reminder":["tinymce"],"remindertranslation":["tinymce"],"reservationitem":["fullcalendar","reservations"]},"management":{"datacenter":{"dcroom":["gridstack","rack"]}},"config":{"commondropdown":{"ITILFollowupTemplate":["tinymce"],"ProjectTaskTemplate":["tinymce"],"SolutionTemplate":["tinymce"],"TaskTemplate":["tinymce"]},"notification":{"notificationtemplate":["tinymce"]},"plugin":{"marketplace":["marketplace"]},"config":["clipboard"]},"admin":{"0":"clipboard","1":"sortable","pluginformcreatorform":["gridstack"]},"preference":["clipboard"],"self-service":["tinymce","fullcalendar","reservations"],"tickets":{"ticket":["tinymce"]},"create_ticket":["tinymce"],"reservation":["tinymce","fullcalendar","reservations"],"faq":["tinymce"]},"max_time_for_count":200,"default_impact_asset_types":{"Appliance":"pics\/impact\/appliance.png","Cluster":"pics\/impact\/cluster.png","Computer":"pics\/impact\/computer.png","Datacenter":"pics\/impact\/datacenter.png","DCRoom":"pics\/impact\/dcroom.png","Domain":"pics\/impact\/domain.png","Enclosure":"pics\/impact\/enclosure.png","Monitor":"pics\/impact\/monitor.png","NetworkEquipment":"pics\/impact\/networkequipment.png","PDU":"pics\/impact\/pdu.png","Peripheral":"pics\/impact\/peripheral.png","Phone":"pics\/impact\/phone.png","Printer":"pics\/impact\/printer.png","Rack":"pics\/impact\/rack.png","Software":"pics\/impact\/software.png","DatabaseInstance":"pics\/impact\/databaseinstance.png"},"impact_asset_types":{"Appliance":"pics\/impact\/appliance.png","Cluster":"pics\/impact\/cluster.png","Computer":"pics\/impact\/computer.png","Datacenter":"pics\/impact\/datacenter.png","DCRoom":"pics\/impact\/dcroom.png","Domain":"pics\/impact\/domain.png","Enclosure":"pics\/impact\/enclosure.png","Monitor":"pics\/impact\/monitor.png","NetworkEquipment":"pics\/impact\/networkequipment.png","PDU":"pics\/impact\/pdu.png","Peripheral":"pics\/impact\/peripheral.png","Phone":"pics\/impact\/phone.png","Printer":"pics\/impact\/printer.png","Rack":"pics\/impact\/rack.png","Software":"pics\/impact\/software.png","DatabaseInstance":"pics\/impact\/databaseinstance.png","AuthLDAP":"pics\/impact\/authldap.png","CartridgeItem":"pics\/impact\/cartridgeitem.png","Contract":"pics\/impact\/contract.png","CronTask":"pics\/impact\/crontask.png","DeviceSimcard":"pics\/impact\/devicesimcard.png","Entity":"pics\/impact\/entity.png","Group":"pics\/impact\/group.png","ITILCategory":"pics\/impact\/itilcategory.png","Line":"pics\/impact\/line.png","Location":"pics\/impact\/location.png","MailCollector":"pics\/impact\/mailcollector.png","Notification":"pics\/impact\/notification.png","Profile":"pics\/impact\/profile.png","Project":"pics\/impact\/project.png","SLM":"pics\/impact\/slm.png","SoftwareLicense":"pics\/impact\/softwarelicense.png","Supplier":"pics\/impact\/supplier.png","User":"pics\/impact\/user.png","Database":"pics\/impact\/database.png"},"root_doc":"","typedoc_icon_dir":"\/pics\/icones","version":"10.0.5","show_jobs_at_login":"1","cut":"360","list_limit":"50","list_limit_max":"60","url_maxlength":"30","event_loglevel":"5","notifications_mailing":"1","admin_email_name":"HELPDESK","from_email":"helpdesk@email.lojasmm.com.br","from_email_name":"HELPDESK","noreply_email":"helpdesk@email.lojasmm.com.br","noreply_email_name":"helpdesk@email.lojasmm.com.br","replyto_email_name":"HELPDESK","mailing_signature":"Boas Vendas!","use_anonymous_helpdesk":"0","use_anonymous_followups":"0","language":"en_US","priority_1":"#ffefdf","priority_2":"#fdd1a9","priority_3":"#ffa05a","priority_4":"#ff4500","priority_5":"#ff2706","priority_6":"#ff0000","date_tax":"2005-12-31","cas_host":"","cas_port":"443","cas_uri":"","cas_logout":"","cas_version":"CAS_VERSION_2_0","existing_auth_server_field_clean_domain":"0","planning_begin":"08:00","planning_end":"20:00","utf8_conv":"1","use_public_faq":"0","url_base":"https:\/\/nexus.lojasmm.com.br","show_link_in_mail":"0","text_login":"","founded_new_version":"","dropdown_max":"100","ajax_wildcard":"*","ajax_limit_count":"10","is_users_auto_add":"1","date_format":"1","csv_delimiter":";","is_ids_visible":"0","smtp_mode":"3","smtp_host":"172.20.61.15","smtp_port":"587","smtp_username":"helpdesk@email.lojasmm.com.br","proxy_name":"","proxy_port":"8080","proxy_user":"","add_followup_on_update_ticket":"1","keep_tickets_on_delete":"1","time_step":"5","helpdesk_doc_url":"","central_doc_url":"","documentcategories_id_forticket":"0","monitors_management_restrict":"2","phones_management_restrict":"2","peripherals_management_restrict":"2","printers_management_restrict":"2","use_log_in_files":"1","time_offset":"10800","is_contact_autoupdate":"1","is_user_autoupdate":"1","is_group_autoupdate":"1","is_location_autoupdate":"1","state_autoupdate_mode":"-1","is_contact_autoclean":"0","is_user_autoclean":"0","is_group_autoclean":"0","is_location_autoclean":"0","state_autoclean_mode":"2","use_flat_dropdowntree":"1","use_autoname_by_entity":"1","softwarecategories_id_ondelete":"1","x509_email_field":"","x509_cn_restrict":"","x509_o_restrict":"","x509_ou_restrict":"","default_mailcollector_filesize_max":"34603008","followup_private":"0","task_private":"0","default_software_helpdesk_visible":"1","names_format":"1","default_requesttypes_id":"1","use_noright_users_add":"1","cron_limit":"5","priority_matrix":{"1":{"1":"1","2":"1","3":"2","4":"2","5":"1"},"2":{"1":"1","2":"2","3":"2","4":"3","5":"2"},"3":{"1":"2","2":"2","3":"3","4":"4","5":"3"},"4":{"1":"2","2":"3","3":"4","4":"4","5":"5"},"5":{"1":"2","2":"3","3":"4","4":"5","5":"5"}},"urgency_mask":"62","impact_mask":"62","user_deleted_ldap":"0","user_restored_ldap":"0","auto_create_infocoms":"1","use_slave_for_search":"0","show_count_on_tabs":"1","refresh_views":"1","set_default_tech":"0","allow_search_view":"2","allow_search_all":"0","allow_search_global":"1","display_count_on_home":"10","use_password_security":"0","password_min_length":"8","password_need_number":"1","password_need_letter":"1","password_need_caps":"1","password_need_symbol":"1","use_check_pref":"0","notification_to_myself":"1","duedateok_color":"#06ff00","duedatewarning_color":"#ffb800","duedatecritical_color":"#ff0000","duedatewarning_less":"60","duedatecritical_less":"10","duedatewarning_unit":"%","duedatecritical_unit":"%","realname_ssofield":"","firstname_ssofield":"","email1_ssofield":"","email2_ssofield":"","email3_ssofield":"","email4_ssofield":"","phone_ssofield":"","phone2_ssofield":"","mobile_ssofield":"","comment_ssofield":"","title_ssofield":"","category_ssofield":"","language_ssofield":"","entity_ssofield":"","registration_number_ssofield":"","ssovariables_id":"0","ssologout_url":"","translate_kb":"0","translate_dropdowns":"0","translate_reminders":"0","pdffont":"helvetica","keep_devices_when_purging_item":"0","maintenance_mode":"0","maintenance_text":"","attach_ticket_documents_to_mail":"1","backcreated":"0","task_state":"1","palette":"auror","page_layout":"vertical","fold_menu":"0","fold_search":"0","savedsearches_pinned":"0","timeline_order":"natural","itil_layout":"","richtext_layout":"classic","lock_use_lock_item":"0","lock_autolock_mode":"1","lock_directunlock_notification":"0","lock_item_list":[],"lock_lockprofile_id":"13","set_default_requester":"1","highcontrast_css":"0","default_central_tab":"0","smtp_check_certificate":"0","enable_api":"1","enable_api_login_credentials":"0","enable_api_login_external_token":"1","url_base_api":"https:\/\/nexus.lojasmm.com.br\/apirest.php\/","login_remember_time":"86400","login_remember_default":"1","use_notifications":"1","notifications_ajax":"0","notifications_ajax_check_interval":"5","notifications_ajax_sound":"sound_a","notifications_ajax_icon_url":"\/pics\/logos\/logo-GLPI-500-black.png","dbversion":"10.0.5@628dbfbb91eb4caf10c35969d9162b9300b141e0","smtp_max_retries":"2","smtp_sender":"helpdesk@email.lojasmm.com.br","instance_uuid":"JMrgUQlhsGIuhwLSJGuAJegmSDPZcVuhdSTvrU0i","registration_uuid":"SPM2FOKOyBhacmu6ihvRkyvSIQJKO4RPG0b4glx7","smtp_retry_time":"1","purge_addrelation":"0","purge_deleterelation":"0","purge_createitem":"0","purge_deleteitem":"0","purge_restoreitem":"0","purge_updateitem":"0","purge_item_software_install":"0","purge_software_item_install":"0","purge_software_version_install":"0","purge_infocom_creation":"0","purge_profile_user":"0","purge_group_user":"0","purge_adddevice":"0","purge_updatedevice":"0","purge_deletedevice":"0","purge_connectdevice":"0","purge_disconnectdevice":"0","purge_userdeletedfromldap":"0","purge_comments":"0","purge_datemod":"0","purge_all":"0","purge_user_auth_changes":"0","purge_plugins":"0","purge_refusedequipment":"0","display_login_source":"1","devices_in_menu":["Item_DeviceSimcard"],"password_expiration_delay":"-1","password_expiration_notice":"-1","password_expiration_lock_delay":"-1","default_dashboard_central":"central","default_dashboard_assets":"assets","default_dashboard_helpdesk":"assistance","default_dashboard_mini_ticket":"mini_tickets","impact_enabled_itemtypes":"[\"Appliance\",\"Cluster\",\"Computer\",\"Datacenter\",\"DCRoom\",\"Domain\",\"Enclosure\",\"Monitor\",\"NetworkEquipment\",\"PDU\",\"Peripheral\",\"Phone\",\"Printer\",\"Rack\",\"Software\",\"DatabaseInstance\"]","document_max_size":"30","planning_work_days":["1","2","3","4","5","6","0"],"system_user":"6","support_legacy_data":"0","marketplace_replace_plugins":"2","glpi_network_uuid":"CaEpxXELlm8SAyvE9e1ZfrfIRTbbErfX3iESPFYt","_matrix":"1","_impact_5":"1","_impact_4":"1","_impact_3":"1","_impact_2":"1","_impact_1":"1","_urgency_5":"1","_matrix_5_5":"5","_matrix_5_4":"5","_matrix_5_3":"4","_matrix_5_2":"3","_matrix_5_1":"2","_urgency_4":"1","_matrix_4_5":"5","_matrix_4_4":"4","_matrix_4_3":"4","_matrix_4_2":"3","_matrix_4_1":"2","_urgency_3":"1","_matrix_3_5":"3","_matrix_3_4":"4","_matrix_3_3":"3","_matrix_3_2":"2","_matrix_3_1":"2","_urgency_2":"1","_matrix_2_5":"2","_matrix_2_4":"3","_matrix_2_3":"2","_matrix_2_2":"2","_matrix_2_1":"1","_urgency_1":"1","_matrix_1_5":"1","_matrix_1_4":"2","_matrix_1_3":"2","_matrix_1_2":"1","_matrix_1_1":"1","_update_devices_in_menu":"1","notification_uuid":"FA6YZJnpPGy1jwaKDpPE66iydmIwO2PeHYGKa8dg","notifications_webhook":"0","update_auth":"Salvar","glpitablesitemtype":{"Plugin":"glpi_plugins","APIClient":"glpi_apiclients","Entity":"glpi_entities"},"glpiitemtypetables":{"glpi_plugins":"Plugin","glpi_apiclients":"APIClient","glpi_entities":"Entity"},"notifications_modes":{"mailing":{"label":"Email","from":"core"},"ajax":{"label":"Browser","from":"core"},"webhook":{"label":"Webhook","from":"webhook"}}};
         var GLPI_PLUGINS_PATH = {"formcreator":"marketplace\/formcreator","behaviors":"marketplace\/behaviors","fields":"marketplace\/fields","addressing":"marketplace\/addressing","moreticket":"marketplace\/moreticket","webhook":"marketplace\/webhook"};


//]]>
</script>
                      </div>
   <div class="floating-buttons d-inline-flex">
      <span class="btn btn-secondary d-none me-1 d-md-block" id="backtotop">
         <i class="fas fa-arrow-up" title="Back to top of the page">
            <span class="visually-hidden">Top of the page</span>
         </i>
      </span>
   </div>

