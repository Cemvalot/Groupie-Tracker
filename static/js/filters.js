// ========================
// Element References
// ========================

// Creation Date Slider
const creationMinRange = document.getElementById("creationMinRange");
const creationMaxRange = document.getElementById("creationMaxRange");
const creationCurrentRange = document.getElementById("creationCurrentRange");
const creationTrack = document.getElementById("creationTrack");

// Album Date Slider
const albumMinRange = document.getElementById("albumMinRange");
const albumMaxRange = document.getElementById("albumMaxRange");
const albumCurrentRange = document.getElementById("albumCurrentRange");
const albumTrack = document.getElementById("albumTrack");

// Location Input & List
const inputField = document.getElementById("inputField");
const addButton = document.getElementById("addButton");
const itemList = document.getElementById("itemList");

// Filter Modal & Button
const filterButton = document.getElementById("filter-button");
const filterModal = document.getElementById("filter-modal");
const closeButton = document.getElementById("closeButton");
const artistGrid = document.getElementById("grid");
const submitButton = document.getElementById("submitButton");
const resetButton = document.getElementById("resetButton");
const searchContainer = document.getElementById("search-container");
const title = document.getElementById("title");

// Performance Locations Array
let performanceLocations = [];

// ========================
// Filter Modal Toggle
// ========================

filterButton.addEventListener("click", (event) => {
  event.stopPropagation(); // Prevent document click from closing the modal
  filterModal.classList.toggle("visible");
  artistGrid.classList.toggle("shifted");
  searchContainer.classList.toggle("shifted");
  title.classList.toggle("shifted");

  if (window.innerWidth <= 768) {
    document.body.style.overflow = "hidden"; // Prevent scrolling on mobile
  }
});

closeButton.addEventListener("click", () => {
  filterModal.classList.remove("visible");
  artistGrid.classList.remove("shifted");
  searchContainer.classList.remove("shifted");
  title.classList.remove("shifted");
  document.body.style.overflow = "auto"; // Restore scrolling
});

// Close sidebar when clicking outside of it
document.addEventListener("click", (event) => {
  if (
    !filterModal.contains(event.target) &&
    event.target !== filterButton &&
    filterModal.classList.contains("visible")
  ) {
    filterModal.classList.remove("visible");
    artistGrid.classList.remove("shifted");
    searchContainer.classList.remove("shifted");
    title.classList.remove("shifted");
  }
});

// ========================
// Creation Date Slider Functions
// ========================

const updateCreationRange = () => {
  const minValue = parseInt(creationMinRange.value);
  const maxValue = parseInt(creationMaxRange.value);

  if (minValue >= maxValue) creationMinRange.value = maxValue - 1;
  if (maxValue <= minValue) creationMaxRange.value = minValue + 1;

  creationCurrentRange.textContent = `${creationMinRange.value} - ${creationMaxRange.value}`;
  updateCreationTrack();
  appendQueryParams();
};

const updateCreationTrack = () => {
  const minValue =
    ((creationMinRange.value - creationMinRange.min) /
      (creationMinRange.max - creationMinRange.min)) *
    100;
  const maxValue =
    ((creationMaxRange.value - creationMinRange.min) /
      (creationMinRange.max - creationMinRange.min)) *
    100;

  creationTrack.style.left = `${minValue}%`;
  creationTrack.style.width = `${maxValue - minValue}%`;
};

// ========================
// Album Date Slider Functions
// ========================

const updateAlbumRange = () => {
  const minValue = parseInt(albumMinRange.value);
  const maxValue = parseInt(albumMaxRange.value);

  if (minValue >= maxValue) albumMinRange.value = maxValue - 1;
  if (maxValue <= minValue) albumMaxRange.value = minValue + 1;

  albumCurrentRange.textContent = `${albumMinRange.value} - ${albumMaxRange.value}`;
  updateAlbumTrack();
  appendQueryParams();
};

const updateAlbumTrack = () => {
  const minValue =
    ((albumMinRange.value - albumMinRange.min) /
      (albumMinRange.max - albumMinRange.min)) *
    100;
  const maxValue =
    ((albumMaxRange.value - albumMinRange.min) /
      (albumMinRange.max - albumMinRange.min)) *
    100;

  albumTrack.style.left = `${minValue}%`;
  albumTrack.style.width = `${maxValue - minValue}%`;
};

// ========================
// Location Functions
// ========================

const formatLocation = (input) => {
  // Trim whitespace from the input
  const trimmedInput = input.trim();

  // Check if input is empty
  if (!trimmedInput) {
    return null;
  }

  // Split by comma to check for "State/City, Country" format
  const parts = trimmedInput.split(",").map((part) => part.trim());

  if (parts.length === 2) {
    const [cityState, country] = parts;

    // Ensure both city/state and country are non-empty
    if (cityState && country) {
      return `${cityState}, ${country}`;
    }
    return null;
  }

  // If no comma, assume it's either "State/City" or "Country"
  return trimmedInput; // return as-is if it doesn't contain a comma
};

// Function to add location to the list
const addLocation = () => {
  const input = inputField.value.trim();
  const formattedLocation = formatLocation(input);

  if (!performanceLocations.includes(formattedLocation)) {
    performanceLocations.push(formattedLocation);
    updateLocationsList();
    inputField.value = "";
    inputField.placeholder = "Add Location...";
  } else {
    inputField.value = "";
    inputField.placeholder = "Location already in list";
    setTimeout(() => {
      inputField.placeholder = "Add Location...";
    }, 2000);
  }
};

// Function to remove location from the list
const removeLocation = (index) => {
  performanceLocations.splice(index, 1);
  updateLocationsList();
};

// Function to update the locations list in the UI
const updateLocationsList = () => {
  itemList.innerHTML = "";
  performanceLocations.forEach((location, index) => {
    const li = document.createElement("li");
    li.textContent = location;
    const removeButton = document.createElement("button");
    removeButton.textContent = "×";
    removeButton.className = "remove-button";
    removeButton.onclick = (event) => {
      event.stopPropagation(); // Stop event bubbling
      removeLocation(index);
    };
    li.appendChild(removeButton);
    itemList.appendChild(li);
  });
};

// ========================
// Query Parameter Management
// ========================

const appendQueryParams = () => {
  const creationMin = creationMinRange.value;
  const creationMax = creationMaxRange.value;
  const albumMin = albumMinRange.value;
  const albumMax = albumMaxRange.value;

  // Get selected band member values
  const selectedBandMembers = Array.from(
    document.querySelectorAll('input[name="bandMembers"]:checked')
  ).map((checkbox) => checkbox.value);

  const params = new URLSearchParams();
  params.set("creationMin", creationMin);
  params.set("creationMax", creationMax);
  params.set("albumMin", albumMin);
  params.set("albumMax", albumMax);

  // Add each selected band member as a separate query parameter
  params.delete("bandMembers"); // Clear any existing bandMembers values
  selectedBandMembers.forEach((value) => {
    params.append("bandMembers", value);
  });

  // Add locations to query parameters
  if (performanceLocations.length > 0) {
    performanceLocations.forEach((location) => {
      params.append("locations", location);
    });
  }

  return params.toString();
};

// ========================
// Event Listeners
// ========================

addButton.addEventListener("click", addLocation);
inputField.addEventListener("keypress", (e) => {
  if (e.key === "Enter") {
    e.preventDefault();
    addLocation();
  }
});

// Initialize from URL parameters
const urlParams = new URLSearchParams(window.location.search);
const locationParams = urlParams.getAll("locations");
if (locationParams.length > 0) {
  performanceLocations = locationParams;
  updateLocationsList();
}

// Updated submit button click handler with forced reload
submitButton.addEventListener("click", () => {
  const queryString = appendQueryParams();
  const baseUrl = window.location.href.split("?")[0];
  filterModal.classList.add("hidden");

  // Force a page reload with the new parameters
  window.location.assign(`${baseUrl}?${queryString}`);
});

// Attach event listeners for creation date slider
creationMinRange.addEventListener("input", updateCreationRange);
creationMaxRange.addEventListener("input", updateCreationRange);

// Attach event listeners for album date slider
albumMinRange.addEventListener("input", updateAlbumRange);
albumMaxRange.addEventListener("input", updateAlbumRange);

// ========================
// Reset Functions
// ========================

const resetCreationRange = () => {
  creationMinRange.value = "1950";
  creationMaxRange.value = "2020";
  updateCreationRange();
};

// Function to reset album date range
const resetAlbumRange = () => {
  albumMinRange.value = "1950";
  albumMaxRange.value = "2020";
  updateAlbumRange();
};

// Function to reset band member checkboxes
const resetBandMembers = () => {
  const checkboxes = document.querySelectorAll('input[name="bandMembers"]');
  checkboxes.forEach((checkbox) => {
    checkbox.checked = false;
  });
};

// Function to reset performance locations
const resetLocations = () => {
  performanceLocations = [];
  updateLocationsList();
  inputField.value = "";
  inputField.placeholder = "e.g. Texas, USA";
};

// Main reset function that calls all individual reset functions
const resetAllFilters = () => {
  resetCreationRange();
  resetAlbumRange();
  resetBandMembers();
  resetLocations();

  // Update the URL to remove all parameters
  const baseUrl = window.location.href.split("?")[0];
  window.history.replaceState({}, "", baseUrl);
};

resetButton.addEventListener("click", resetAllFilters);

// ========================
// Initialization
// ========================

updateCreationTrack();
updateAlbumTrack();

//
