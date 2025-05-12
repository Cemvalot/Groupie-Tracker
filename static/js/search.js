// Configuration Constants
const MIN_INPUT_LENGTH = 2; // Minimum characters required to fetch suggestions
const SUGGESTIONS_API = "/search"; // API endpoint for fetching suggestions

// Main Function to Handle Suggestions
async function fetchSuggestions() {
  const searchInput = document.getElementById("search-input");
  const suggestionsBox = document.getElementById("suggestions");
  const searchType = document.querySelector("select[name='searchType']").value;
  const query = searchInput.value.trim();

  // Clear previous suggestions
  clearSuggestions(suggestionsBox);

  // Only fetch suggestions if the query length meets the threshold
  if (query.length < MIN_INPUT_LENGTH) return;

  try {
    const suggestions = await fetchSuggestionsFromServer(query, searchType);
    renderSuggestions(suggestions, suggestionsBox);
  } catch (error) {
    console.error("Error fetching suggestions:", error);
    renderErrorMessage(
      suggestionsBox,
      "Failed to fetch suggestions. Please try again."
    );
  }
}

// Fetch Suggestions from Backend
async function fetchSuggestionsFromServer(query, searchType) {
  const response = await fetch(
    `${SUGGESTIONS_API}?q=${encodeURIComponent(
      query
    )}&searchType=${encodeURIComponent(searchType)}`
  );
  if (!response.ok) {
    throw new Error("Network response was not ok");
  }
  return response.json();
}

// Render Suggestions in the Dropdown
function renderSuggestions(suggestions, suggestionsBox) {
  suggestionsBox.innerHTML = suggestions
    .map(
      (suggestion) => `
        <div class="suggestion-item" data-input="${suggestion.input}">
            ${suggestion.value}
        </div>
    `
    )
    .join("");
  suggestionsBox.classList.add("show");
}

// Clear Previous Suggestions
function clearSuggestions(suggestionsBox) {
  suggestionsBox.innerHTML = "";
  suggestionsBox.classList.remove("show");
}

// Render Error Message
function renderErrorMessage(suggestionsBox, message) {
  suggestionsBox.innerHTML = `<div class="error-message">${message}</div>`;
  suggestionsBox.classList.add("show");
}

// Event Listener for Selecting a Suggestion
function setupSuggestionSelection() {
  const searchInput = document.getElementById("search-input");
  const suggestionsBox = document.getElementById("suggestions");

  suggestionsBox.addEventListener("click", (event) => {
    if (event.target.classList.contains("suggestion-item")) {
      let input = event.target.getAttribute("data-input");

      // Remove any extra context from the suggestion (e.g., " - member of")
      if (input.includes(" - ")) {
        input = input.split(" - ")[0];
      }

      // Update search input value
      searchInput.value = input;

      // Log query for debugging
      console.log(`Searching for: ${input}`);

      // Automatically trigger form submission
      const searchForm = document.querySelector("form.search-form");
      if (searchForm) {
        searchForm.submit();
      }

      clearSuggestions(suggestionsBox);
    }
  });

  document.addEventListener("click", (event) => {
    if (!event.target.closest(".search-container")) {
      clearSuggestions(suggestionsBox);
    }
  });
}

// Initialize Event Listeners
window.addEventListener("DOMContentLoaded", () => {
  setupSuggestionSelection();
});

// Add "General Search" option in dropdown initialization
const searchTypeDropdown = document.querySelector("select[name='searchType']");
if (searchTypeDropdown) {
  const generalOption = document.createElement("option");
  generalOption.value = "general";
  generalOption.textContent = "General Search";
  generalOption.defaultSelected = true;
  searchTypeDropdown.prepend(generalOption);
}
