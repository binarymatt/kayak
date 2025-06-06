export async function kayak<T>(url: string): Promise<T> {
  console.log("calling kayak");
  const response = await fetch(url, {
    method: "POST",
    mode: "no-cors",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(response.statusText);
  }

  //    And can also be used here ↴
  return (await response.json()) as T;
}
