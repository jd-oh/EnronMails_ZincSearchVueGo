// language: javascript
export async function fetchFolders() {
    try {
      const response = await fetch('http://localhost:8080/api/folders');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      return await response.json();
    } catch (error) {
      console.error("Error fetching folders:", error);
      throw error;
    }
  }