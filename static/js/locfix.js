document.addEventListener("DOMContentLoaded", function () {
  // Helper function to convert a string to title case
  function toTitleCase(str) {
    return str
      .toLowerCase()
      .split(" ")
      .map(function (word) {
        return word.charAt(0).toUpperCase() + word.slice(1);
      })
      .join(" ");
  }

  // Clean and title-case location names
  const locationElems = document.querySelectorAll(".location-name");
  locationElems.forEach((elem) => {
    let text = elem.textContent;
    // Replace underscores and hyphens with a space
    text = text.replace(/[_-]/g, " ");
    // Convert to title case
    text = toTitleCase(text);
    elem.textContent = text;
  });
});

document.addEventListener("DOMContentLoaded", function () {
  const backContents = document.querySelectorAll(".back-content");

  backContents.forEach((backContent) => {
    // Grab all the <p class="date"> elements inside this card
    const dateElements = backContent.querySelectorAll(".concert-dates .date");

    // If there's only one date, center it vertically & horizontally
    if (dateElements.length <= 8) {
      backContent.style.justifyContent = "center";
      backContent.style.alignItems = "center";
      backContent.style.overflowY = "hidden"; // remove scrollbar if you want
    }
  });
});

document.addEventListener("DOMContentLoaded", function () {
  const locationCards = document.querySelectorAll(".location-card");
  const detailsBox = document.querySelector(".details-box");

  // For example, if there are more than 5 location cards, enlarge the box
  if (locationCards.length > 5) {
    detailsBox.classList.add("bigger-details-box");
  }
});
