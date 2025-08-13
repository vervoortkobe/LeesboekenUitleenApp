function zoekFunction() {
  let input = document.getElementById("zoek").value.toLowerCase();
  let items = document.querySelectorAll(".card");
  items.forEach((item) => {
    item.style.display = item.textContent.toLowerCase().includes(input)
      ? ""
      : "none";
  });
}

function voegKlasToe() {
  let naam = prompt("Klasnaam:");
  if (naam) {
    fetch("/klassen", {
      method: "POST",
      body: JSON.stringify({ naam }),
      headers: { "Content-Type": "application/json" },
    }).then(() => location.reload());
  }
}

function pasKlasAan(id) {
  let naam = prompt("Nieuwe klasnaam:");
  if (naam) {
    fetch(`/klassen/${id}`, {
      method: "PUT",
      body: JSON.stringify({ naam }),
      headers: { "Content-Type": "application/json" },
    }).then(() => location.reload());
  }
}

function verwijderKlas(id) {
  if (confirm("Klas verwijderen?")) {
    fetch(`/klassen/${id}`, { method: "DELETE" }).then(() => location.reload());
  }
}

function voegLeerlingToe(klasID) {
  let naam = prompt("Leerlingnaam:");
  if (naam) {
    fetch("/leerlingen", {
      method: "POST",
      body: JSON.stringify({ naam, klasID }),
      headers: { "Content-Type": "application/json" },
    }).then(() => location.reload());
  }
}

function pasLeerlingAan(id) {
  let naam = prompt("Nieuwe leerlingnaam:");
  if (naam) {
    fetch(`/leerlingen/${id}`, {
      method: "PUT",
      body: JSON.stringify({ naam }),
      headers: { "Content-Type": "application/json" },
    }).then(() => location.reload());
  }
}

function verwijderLeerling(id) {
  if (confirm("Leerling verwijderen?")) {
    fetch(`/leerlingen/${id}`, { method: "DELETE" }).then(() =>
      location.reload()
    );
  }
}

function voegNiveauToe() {
  let niveau = prompt("AVI-niveau (bijv. AVI1):");
  if (niveau) {
    fetch("/niveaus", {
      method: "POST",
      body: JSON.stringify({ niveau }),
      headers: { "Content-Type": "application/json" },
    }).then(() => location.reload());
  }
}

function pasNiveauAan(id) {
  let niveau = prompt("Nieuw AVI-niveau:");
  if (niveau) {
    fetch(`/niveaus/${id}`, {
      method: "PUT",
      body: JSON.stringify({ niveau }),
      headers: { "Content-Type": "application/json" },
    }).then(() => location.reload());
  }
}

function voegBoekToe(niveauID) {
  let titel = prompt("Boektitel:");
  if (titel) {
    fetch("/boeken", {
      method: "POST",
      body: JSON.stringify({ titel, aviNiveauID: niveauID }),
      headers: { "Content-Type": "application/json" },
    }).then(() => location.reload());
  }
}

function pasBoekAan(id) {
  let titel = prompt("Nieuwe boektitel:");
  if (titel) {
    fetch(`/boeken/${id}`, {
      method: "PUT",
      body: JSON.stringify({ titel }),
      headers: { "Content-Type": "application/json" },
    }).then(() => location.reload());
  }
}

function verwijderBoek(id) {
  if (confirm("Boek verwijderen?")) {
    fetch(`/boeken/${id}`, { method: "DELETE" }).then(() => location.reload());
  }
}

function updateLeerlingNaam(id) {
  let naam = document.getElementById("leerlingNaam").value;
  fetch(`/leerlingen/${id}/naam`, {
    method: "PUT",
    body: JSON.stringify({ naam }),
    headers: { "Content-Type": "application/json" },
  }).then(() => location.reload());
}

function updateLeesDatum(boekID, leerlingID, datum) {
  fetch(`/leerlingen/${leerlingID}/leesdata`, {
    method: "POST",
    body: JSON.stringify({ boekID, datum }),
    headers: { "Content-Type": "application/json" },
  }).then(() => location.reload());
}
