import React, { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom';
import './Home.css'

// Product Card Component
function ProductCard({ product }) {
  return (
    <div className="card">
      <img src={`${process.env.REACT_APP_IMAGE_STORAGE}\\product_photo\\${product.product_photo}`} alt={product.product_name} />
      <h2>{product.product_name}</h2>
      <p>Price: ${product.price}</p>
      <p>Visibility: {product.visibility ? "public":"private"}</p>
      <p>{product.username? product.username:""}</p>
    </div>
  );
}

function Home() {
  const [publicProducts, setPublicProducts] = useState([]);
  const [myProducts, setMyProducts] = useState([]);
  const navigate = useNavigate()

  useEffect(() => {
    if (myProducts.length === 0 || publicProducts.length === 0){
      fetchProducts();
    }
  }, []); 

  const fetchProducts = async () => {
    try {
      const response = await fetch(`${process.env.REACT_APP_API_HOST}/api/product`,{
        method: "GET",
        headers: {
          'Content-Type': 'application/json',
        },
        credentials:'include'
      });
      const data = await response.json();
      if (data && !data.error) {
        setPublicProducts(data);
      } else {
        throw new Error(data?.error || "Failed to fetch public products.");
      }
    } catch (err) {
      console.error('Error fetching products:', err);
      if (!err.message.includes('Failed to fetch')) {
        navigate("/login");
      }
    }

    try {
      const response2 = await fetch(`${process.env.REACT_APP_API_HOST}/api/product?isPrivate=true`,{
        method: "GET",
        headers: {
          'Content-Type': 'application/json',
        },
        credentials:'include'
      });
      const data2 = await response2.json();
      if (data2 && !data2.error) {
        setMyProducts(data2);
      } else {
        throw new Error(data2?.error || "Failed to fetch private products.");
      }
    } catch (err) {
      console.error('Error fetching products:', err);
      if (!err.message.includes('Failed to fetch')) {
        navigate("/login");
      }
    }
  };

  return (
    <div>
      <div className='title-head'>
      <h1>Home</h1>
      <h1><Link to="/post_product">Post Product</Link></h1>
      </div>
      <div className='product-box'>
        <div className='product-display'>
          <h3>Public Product</h3>
          <div className="product-grid">
            {publicProducts.map((product) => (
              <ProductCard key={product.product_id} product={product} />
            ))}
          </div>
        </div>
        <div className='product-display'>
        <h3>My Product</h3>
        <div className="product-grid">
          {myProducts.map((product) => (
            <ProductCard key={product.product_id} product={product} />
          ))}
        </div>
        </div>
      </div>
    </div>
  );
}

export default Home