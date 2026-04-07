import { useEffect, useState } from "react";

type Language = {
  name: string;
  votes: number;
  image: string;
};

type WindowWithRuntimeConfig = Window & {
  __RUNTIME_CONFIG__?: {
    VITE_API_URL?: string;
  };
};

// Read API URL from environment (Vite exposes variables prefixed with VITE_)
// const API_URL = import.meta.env.VITE_API_URL || "http://127.0.0.1:3000";
const API_URL =
  (window as WindowWithRuntimeConfig).__RUNTIME_CONFIG__?.VITE_API_URL

function App() {
  const [languages, setLanguages] = useState<Language[]>([]);

  // Fetch languages from Go API
  useEffect(() => {
    fetch(`${API_URL}/languages`)
      .then((res) => res.json())
      .then((data) => setLanguages(data))
      .catch((err) => console.error("Error fetching languages:", err));
  }, []);

  // Vote handler (POST /languages)
  const vote = async (name: string, delta: number) => {
    try {
      const res = await fetch(`${API_URL}/languages`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, delta }),
      });
      const updated = await res.json();

      // Update local state with returned language object
      setLanguages(
        languages.map((l) => (l.name === updated.name ? updated : l)),
      );
    } catch (err) {
      console.error("Error voting:", err);
    }
  };

  return (
    <div className="min-h-screen flex flex-col bg-cyan-700">
      {/* Navigation Bar */}
      <nav className="bg-yellow-600 text-white shadow-md">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-center">
          <h1 className="text-2xl font-bold">Voting App</h1>
        </div>
      </nav>

      {/* Title */}
      <header className="text-center mt-8 mb-6">
        <h2 className="text-3xl font-semibold text-white">
          Vote Your Favourite Languages
        </h2>
      </header>

      {/* Grid of Cards */}
      <main className="grow">
        <div className="max-w-7xl mx-auto px-4 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {languages.map((lang) => (
            <div
              key={lang.name}
              className="bg-yellow-300/40 backdrop-blur-md rounded-xl shadow-lg overflow-hidden flex flex-col items-center p-6 hover:scale-105 transform transition"
            >
              <img
                src={lang.image}
                alt={lang.name}
                className="w-24 h-24 object-contain mb-4"
              />
              <h2 className="text-xl text-white font-bold mb-2">{lang.name}</h2>
              <p className="text-white font-bold mb-4">Votes: {lang.votes}</p>
              <div className="flex gap-4">
                <button
                  onClick={() => vote(lang.name, 1)}
                  className="px-4 py-2 bg-green-500 font-bold text-white rounded hover:bg-green-600 cursor-pointer"
                >
                  +
                </button>
                <button
                  onClick={() => vote(lang.name, -1)}
                  className="px-4 py-2 font-bold bg-red-500 text-white rounded hover:bg-red-600 cursor-pointer"
                >
                  -
                </button>
              </div>
            </div>
          ))}
        </div>
      </main>

      {/* Footer */}
      <footer className="bg-yellow-600 text-white font-bold mt-8">
        <div className="max-w-7xl mx-auto px-4 py-4 text-center">
          <p className="text-sm">
            &copy; 2026 Voting App. All rights reserved.
          </p>
        </div>
      </footer>
    </div>
  );
}

export default App;
