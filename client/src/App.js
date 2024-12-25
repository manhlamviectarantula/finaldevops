import './App.css';
import axios from 'axios';
import {useEffect, useState} from 'react';

function App() {

  const [products, setProducts] = useState([])

  useEffect(() => {
    getProducts();
  }, []);

  const getProducts = async () => {
    try {
      const response = await axios.get('http://localhost:8080/api/v1/get-products')
      //nên dùng cách này thay vì response.data, khi ko có data sẽ ko bị trắng trang
      const {data = []} = response
      setProducts(data)
      console.log(response.data)
    } catch (error) {
      console.error(error)
    }
  }

  return (
    <div className="App">
      <h1>react</h1>
      {
        products.map(product => (
          (
            <div key={product._id}>
              <h2>{product.name}</h2>
              <p>{product.description}</p>
              <p>{product.price}</p>
            </div>
          )
        ))
      }
    </div>
  );
}

export default App;
