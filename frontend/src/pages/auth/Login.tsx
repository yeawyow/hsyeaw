import { Input, Button } from "@headlessui/react";

export default function Login() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-4 space-y-4 antialiased text-gray-900 bg-gray-100 dark:bg-dark dark:text-light">
      <div className="w-full max-w-sm px-4 py-6 space-y-6 bg-white rounded-md dark:bg-darker ">
        <h1 className="text-xl font-semibold text-center">Login</h1>

        <form className="space-y-6">
          <Input
            name="user"
            type="text"
            className={`w-full px-4 py-2 border border-gray-300 rounded-md dark:bg-darker dark:border-gray-700 focus:outline-none focus:ring focus:ring-primary-100 dark:focus:ring-primary-darker`}
            placeholder="username"
          />
          <Input
            name="password"
            type="password"
            className={`w-full px-4 py-2 border border-gray-300 rounded-md dark:bg-darker dark:border-gray-700 focus:outline-none focus:ring focus:ring-primary-100 dark:focus:ring-primary-darker`}
            placeholder="password"
          />
          <Button type="submit" className={``}>
            Login
          </Button>
        </form>
      </div>
    </div>
  );
}
