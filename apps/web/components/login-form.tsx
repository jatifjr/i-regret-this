import { Label } from "@workspace/ui/components/label";
import { Input } from "@workspace/ui/components/input";
import { Button } from "@workspace/ui/components/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@workspace/ui/components/card";
import { LogIn } from "lucide-react";

export function LoginForm() {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-center text-5xl font-bold">
          UNW TOEFL
        </CardTitle>
      </CardHeader>
      <CardContent>
        <form className="grid items-center w-full gap-4">
          <div className="flex flex-col gap-2">
            <Label htmlFor="username">Nomor Induk Mahasiswa</Label>
            <Input
              id="username"
              type="username"
              placeholder="NIM.123456"
              required
            ></Input>
          </div>
          <div className="flex flex-col gap-2">
            <Label htmlFor="password">Password</Label>
            <Input
              id="password"
              type="password"
              placeholder="Password"
              required
            ></Input>
          </div>
        </form>
      </CardContent>
      <CardFooter>
        <div className="grid items-center w-full gap-4">
          <div className="flex flex-col gap-2">
            <CardDescription
              hidden={false}
              className="text-center text-red-500"
            >
              NIM atau password yang anda masukkan salah
            </CardDescription>
            <Button>
              <h1>Masuk</h1>
              <LogIn />
            </Button>
          </div>
        </div>
      </CardFooter>
    </Card>
  );
}
