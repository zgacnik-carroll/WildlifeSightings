// Confirm before submitting a new sighting if fields are empty
document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("new-sighting-form");
    if (form) {
        form.addEventListener("submit", (e) => {
            const animal = document.querySelector("#animal").value.trim();
            const location = document.querySelector("#location").value.trim();
            if (!animal || !location) {
                e.preventDefault();
                alert("Please fill in both the animal and location fields.");
            }
        });
    }
});