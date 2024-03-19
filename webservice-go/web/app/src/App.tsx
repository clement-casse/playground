
export function App() {
  return (
    <>
      <header className="sticky top-0 bg-white shadow">
        <div className="container flex flex-col sm:flex-row justify-between items-center mx-auto py-2 px-8">
          <div className="flex items-center text-2xl">
            Title
          </div>
          <nav className="flex mt-4 sm:mt-0">
            <NavItems />
          </nav>
          <div className="hidden md:block">
            <button type="button" className="py-3 px-8 text-sm bg-teal-500 hover:bg-teal-600 rounded text-white">
              Some Button
            </button>
          </div>
        </div>
      </header>

      <main className="text-gray-900">

      </main>

      <footer className="container mx-auto py-16 px-3 mt-48 mb-8 text-gray-800">
        <div className="flex -mx-3">
          <div className="flex-1 px-3">
            <h2 className="text-lg font-semibold">About</h2>
            <p></p>
          </div>
          <div className="flex-1 px-3">
            <h2 className="text-lg font-semibold"></h2>
          </div>
          <div className="flex-1 px-3">
            <h2 className="text-lg font-semibold">Links</h2>
            <ul className="mt-4 leading-loose">
              <li><a href="https://github.com/clement-casse/playground/tree/main/webservice-go">GitHub Repo</a></li>
              <li></li>
            </ul>
          </div>
        </div>
      </footer>
    </>
  )
}

export function NavItems() {
  return (
    <ul>
      <li className="px-4"><a href="/">Home</a></li>
    </ul>
  )
}