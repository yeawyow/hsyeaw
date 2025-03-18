import React, { useState, useEffect } from "react";
// import { useLocation } from "react-router-dom";
import Loader from "./common/loader";
import RouterApp from "./route";

export default function App() {
  const [loading, setLoading] = useState<boolean>(false);
  // const { pathname } = useLocation();

  // useEffect(() => {
  //   window.scrollTo(0, 0);
  // }, [pathname]);
  useEffect(() => {
    setTimeout(() => setLoading(false), 1000); // ตั้งค่า loading = false หลัง 1 วิ
  }, []);

  return <div>{loading ? <Loader /> : <RouterApp />}</div>;
}
