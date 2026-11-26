import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { createBrowserRouter, RouterProvider } from "react-router";
import "@fontsource-variable/geist/index.css";
import "@fontsource-variable/geist-mono/index.css";
import "./styles.css";
import { Layout } from "./components/Layout";
import { I18nProvider } from "./i18n";
import { Channel } from "./pages/Channel";
import { Dashboard } from "./pages/Dashboard";
import { Landing } from "./pages/Landing";
import { Imprint, NotFound, Privacy } from "./pages/Legal";

const queryClient = new QueryClient({
  defaultOptions: { queries: { staleTime: 15_000, refetchOnWindowFocus: false } },
});

const router = createBrowserRouter([
  {
    element: <Layout />,
    children: [
      { path: "/", element: <Landing /> },
      { path: "/dashboard", element: <Dashboard /> },
      { path: "/c/:channel", element: <Channel /> },
      // Links of the old site.
      { path: "/leaderboard/:channel", element: <Channel /> },
      { path: "/imprint", element: <Imprint /> },
      { path: "/privacy", element: <Privacy /> },
      { path: "*", element: <NotFound /> },
    ],
  },
]);

const root = document.getElementById("root");
if (!root) throw new Error("missing #root");

createRoot(root).render(
  <StrictMode>
    <I18nProvider>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </I18nProvider>
  </StrictMode>,
);
