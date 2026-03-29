"use client";
import { Amplify } from "aws-amplify";
import { amplifyConfig } from "@/lib/amplifyConfig";

Amplify.configure(amplifyConfig);

export default function AmplifyProvider({ children }) {
  return children;
}
