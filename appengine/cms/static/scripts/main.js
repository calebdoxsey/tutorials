function api(method, endpoint, data, callback) {
  var xhr = new XMLHttpRequest();
  xhr.open(method, endpoint);
  if (data) {
    xhr.send(JSON.stringify(data));
  } else {
    xhr.send(null);
  }
  xhr.onreadystatechange = function(evt) {
    if (xhr.readyState === 4) {
      var res;
      try {
        res = JSON.parse(xhr.responseText);
      } catch(e) {
        res = { "error": xhr.responseText };
      }
      if (res && res.error) {
        callback(null, res.error);
      } else {
        callback(res);
      }
    }
  };
}

function onLogin(evt) {
  evt.preventDefault();

  var email = document.getElementById("email").value;
  var password = document.getElementById("password").value;

  api("POST", "/api/users/login", {
    "email": email,
    "password": password
  }, function() {
    location.href = "/";
  });
}

function onDocumentCreate(evt) {
  evt.preventDefault();

  var link = document.getElementById("link").value;
  var contents = document.getElementById("contents").value;

  api("POST", "/api/documents", {
    "link": link,
    "contents": contents
  }, function(document, error) {
    if (error) {
      alert(error);
      return;
    }
    location.href = "/documents/" + document.ID;
  })
}

function onDocumentDelete(evt) {
  evt.preventDefault();
}

function onDocumentUpdate(evt) {
  evt.preventDefault();

  var id = document.getElementById("document-id").value;
  var link = document.getElementById("document-link").value;
  var contents = document.getElementById("document-contents").value;

  api("PUT", "/api/documents/" + id, {
    "Link": link,
    "Contents": contents
  }, function(document, error) {
    if (error) {
      alert(error);
      return;
    }
    location.reload();
  });
}

function main() {
  var prev = null;
  function renderView(url) {
    if (prev) {
      document.body.removeChild(prev);
    }
    var view = Views[url];
    if (!view) {
      var arr = url.split("/");
      arr.pop();
      view = Views[arr.join("/") + "/"];
    }
    prev = view();
    document.body.appendChild(prev);
  }
  window.addEventListener("hashchange", function() {
    renderView(location.hash.substr(1) || "/");
  }, false);
  renderView(location.hash.substr(1) || "/");


  /*
  var lf = document.getElementById("login-form");
  if (lf) {
    lf.addEventListener("submit", onLogin);
  }

  var dl = document.getElementById("documents-list");
  if (dl) {
    api("GET", "/api/documents", null, function(docs) {
      for (var i=0; i<docs.length; i++) {
        var doc = docs[i];
        dl.appendChild(
          h("li",
            h("a", {
              "href": "/documents/" + doc.ID
            }, doc.Link),
            " ",
            h("a.documents-delete", {
              "href": "#remove",
              "onclick": onDocumentDelete
            }, "x")
          )
        );
      }
    });
  }

  var dc = document.getElementById("documents-create");
  if (dc) {
    dc.addEventListener("submit", onDocumentCreate);
  }

  var du = document.getElementById("documents-update");
  if (du) {
    var id = document.getElementById("document-id").value;
    api("GET", "/api/documents/" + id, null, function(doc, error) {
      if (error) {
        alert("Error getting document: " + error);
        return;
      }
      document.getElementById("document-link").value = doc.Link;
      document.getElementById("document-contents").value = doc.Contents;
    });
    du.addEventListener("submit", onDocumentUpdate);
  }
  */
}

main();
