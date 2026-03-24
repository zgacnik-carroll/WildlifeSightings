// Confirm before submitting a new sighting if fields are empty
document.addEventListener("DOMContentLoaded", () => {
    const form = document.querySelector("form");
    if (form) {
        form.addEventListener("submit", (e) => {
            const animal = form.querySelector("#animal")?.value.trim();
            const location = form.querySelector("#location")?.value.trim();
            if (!animal || !location) {
                e.preventDefault();
                alert("Please fill in both the animal and location fields.");
            }
        });
    }
});