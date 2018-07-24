function PrepareLS() {
  if (localStorage.getItem("paths") === null) {
    var paths = {};
    localStorage.setItem("paths", JSON.stringify(paths));
  }
}

function PutIntoLS(dbPath) {
  var paths = JSON.parse(localStorage.getItem("paths"));
  if (paths[dbPath] == null) {
    paths[dbPath] = { uses: 1 };
  } else {
    paths[dbPath].uses += 1;
  }

  localStorage.setItem("paths", JSON.stringify(paths));
}

function GetPaths() {
  var paths = JSON.parse(localStorage.getItem("paths"));

  // Sorting. Return only keys;
  var sortedPaths = Object.keys(paths).sort(function(a, b) {
    if (paths[a].uses < paths[b].uses) {
      return 1;
    }
    if (paths[a].uses > paths[b].uses) {
      return -1;
    }
    return 0;
  });

  return sortedPaths;
}

function DeletePath(path) {
  var paths = JSON.parse(localStorage.getItem("paths"));
  delete paths[path];
  localStorage.setItem("paths", JSON.stringify(paths));

  ShowPathsForDeleting();
}
